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

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 14 {
			continue
		}

		return &CheckMResult{
			Completeness:  fields[5],
			Contamination: fields[6],
			GenomeSize:    fields[8],
			N50:           fields[13],
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

	if bestName == "" {
		return "", errors.New("no valid data found in fastani result")
	}
	if bestANI < minANI {
		return "", fmt.Errorf("best ANI match below species threshold: %.2f%%",
			bestANI)
	}
	return bestName, nil
}

type MLSTResult struct {
	Scheme string
	ST     string
}

func ParseMLST(filePath string) (*MLSTResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open mlst result: %v", err)
	}
	defer file.Close()

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

		return &MLSTResult{
			Scheme: fields[1],
			ST:     fields[2],
		}, nil
	}

	return nil, errors.New("No valid data found in mlst result")
}
