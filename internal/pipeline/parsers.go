package pipeline

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CheckMResult struct {
	Completeness  string
	Contamination string
	GenomeSize    string
	Contigs       string
	N50           string
}

func ParseCheckM(filePath string) (*CheckMResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open checkm result: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("Error reading checkm file: %v", err)
		}
		return nil, errors.New("Empty checkm result")
	}

	columns := make(map[string]int)
	for i, col := range strings.Split(strings.TrimSpace(scanner.Text()), "\t") {
		columns[col] = i
	}

	for _, key := range []string{"Name", "Completeness", "Contamination",
		"Genome_Size", "Total_Contigs", "Contig_N50"} {
		if _, ok := columns[key]; !ok {
			return nil, fmt.Errorf("Missing column %s in checkm result", key)
		}
	}
	name, completeness := columns["Name"], columns["Completeness"]
	contamination, genomeSize := columns["Contamination"], columns["Genome_Size"]
	contigs, n50 := columns["Total_Contigs"], columns["Contig_N50"]

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) <= n50 || fields[name] == "" {
			continue
		}

		return &CheckMResult{
			Completeness:  fields[completeness],
			Contamination: fields[contamination],
			GenomeSize:    fields[genomeSize],
			Contigs:       fields[contigs],
			N50:           fields[n50],
		}, nil
	}

	return nil, errors.New("No valid data found in checkm result")
}

func ParseFastANI(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Failed to open fastani result: %v", err)
	}
	defer file.Close()

	const minANI = 95.0

	var bestName string
	var bestANI float64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 3 {
			continue
		}
		ani, err := strconv.ParseFloat(fields[2], 64)
		if err != nil || ani <= bestANI {
			continue
		}
		pathParts := strings.Split(fields[1], "/")
		filename := pathParts[len(pathParts)-1]
		nameParts := strings.Split(filename, ".")
		if len(nameParts) > 0 {
			bestANI = ani
			bestName = nameParts[0]
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Error reading fastani result: %v", err)
	}

	if bestName == "" {
		return "", errors.New("no valid data found in fastani result")
	}
	if bestANI < minANI {
		return "", fmt.Errorf("best ANI match below species threshold: %.2f%%",
			bestANI)
	}
	return bestName, nil
}

func ParseMLST(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Failed to open mlst result: %v", err)
	}
	defer file.Close()

	scheme, st := "-", "-"

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) < 3 {
			continue
		}

		scheme = fields[1]
		st = fields[2]

		if scheme != "-" && st != "-" {
			return fmt.Sprintf("%s (ST%s)", scheme, st), nil
		} else if scheme != "-" && st == "-" {
			return fmt.Sprintf("%s (New ST)", scheme), nil
		} else if scheme == "-" && st == "-" {
			return "Not available for this species", nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Error reading mlst result: %v", err)
	}

	return "", errors.New("No valid data found in mlst result")
}
