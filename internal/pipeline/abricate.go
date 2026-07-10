package pipeline

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var vanPattern = regexp.MustCompile(`(?i)^Van`)

func GetAbricateResult(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open Abricate result: %v", err)
	}
	defer file.Close()

	br := bufio.NewReader(file)
	_, err = br.Peek(1)
	if err == io.EOF {
		return nil, errors.New("Empty Abricate result")
	}

	var results []string
	scanner := bufio.NewScanner(br)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Split(line, "\t")

		if len(fields) < 11 {
			continue
		}

		coverage, errCov := strconv.ParseFloat(fields[9], 64)
		identity, errId := strconv.ParseFloat(fields[10], 64)

		if errCov != nil || errId != nil {
			continue
		}

		gene := fields[5]

		if (coverage > 90.0 && identity > 90.0) ||
			vanPattern.MatchString(gene) {
			results = append(results, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf(
			"An error occurred while reading the file: %v", err)
	}

	return results, nil
}

func ProcessResfinder(abricateResult []string, refCatalogPath string) (
	[]string, []string, error) {
	var geneResults []string
	var blastOutResults []string
	var refList [][]string

	refFile, err := os.Open(refCatalogPath)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Failed to open Resfinder reference file:%v", err)
	}
	defer refFile.Close()

	br := bufio.NewReader(refFile)
	_, err = br.Peek(1)
	if err == io.EOF {
		return nil, nil, errors.New("Empty Resfinder reference file")
	}

	refScanner := bufio.NewScanner(refFile)
	for refScanner.Scan() {
		line := strings.TrimSpace(refScanner.Text())
		if line != "" {
			refList = append(refList, strings.Split(line, "\t"))
		}
	}

	for _, line := range abricateResult {
		fields := strings.Split(line, "\t")
		if len(fields) < 11 {
			continue
		}

		gene := fields[5]
		covDb := fields[6]
		covQ := fields[9]
		id := fields[10]

		nameGeneParts := strings.Split(gene, "_")
		baseNameGene := nameGeneParts[0]

		foundAntibiotic := false
		for _, refItem := range refList {
			if len(refItem) < 17 {
				continue
			}

			if strings.Contains(strings.ToLower(refItem[0]),
				strings.ToLower(baseNameGene)) {
				antibioticName := strings.ToLower(refItem[len(refItem)-17])
				geneResults = append(geneResults, fmt.Sprintf(
					"%s (resistance to %s) (allele confidence %s)", gene,
					antibioticName, id))
				foundAntibiotic = true
				break
			}
		}

		if !foundAntibiotic {
			geneResults = append(geneResults, fmt.Sprintf(
				"%s (allele confidence %s)", gene, id))
		}

		blastOut := fmt.Sprintf("%s (ID: %s COV_Q: %s COV_DB: %s)", gene, id,
			covQ, covDb)
		blastOutResults = append(blastOutResults, blastOut)
	}

	return geneResults, blastOutResults, nil
}

func ProcessVFDB(abricateResult []string) []string {
	var results []string

	for _, line := range abricateResult {
		fields := strings.Split(line, "\t")
		if len(fields) < 14 {
			continue
		}

		results = append(results,
			fmt.Sprintf("%s: %s %s ID: %s COV_Q: %s COV_DB: %s| ",
				fields[1], fields[5], fields[13], fields[10], fields[9],
				fields[6]))
	}

	return results
}

func ProcessPlasmidFinder(abricateResult []string) []string {
	var results []string

	for _, line := range abricateResult {
		fields := strings.Split(line, "\t")
		if len(fields) < 11 {
			continue
		}

		results = append(results,
			fmt.Sprintf("%s (ID: %s COV_Q: %s COV_DB: %s)",
				fields[5], fields[10], fields[9], fields[6]))
	}

	return results
}
