package pipeline

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type KrakenSpecies struct {
	Name  string
	Count int
}

func KrakenSpeciesCounter(krakenOutput string) (*KrakenSpecies, *KrakenSpecies,
	error) {
	file, err := os.Open(krakenOutput)
	if err != nil {
		return nil, nil, fmt.Errorf("Kraken output file not found: %v", err)
	}
	defer file.Close()

	br := bufio.NewReader(file)
	_, err = br.Peek(1)
	if err == io.EOF {
		return nil, nil, errors.New("Empty Kraken result")
	}

	counts := make(map[string]int)
	scanner := bufio.NewScanner(br)

	const maxCapacity = 1024 * 1024 * 50
	lineBuf := make([]byte, 1024*64)
	scanner.Buffer(lineBuf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.SplitN(line, "\t", 4)

		if len(fields) >= 3 && fields[0] == "C" {
			speciesPart := fields[2]
			speciesName := strings.TrimSpace(strings.Split(
				speciesPart, "(")[0])

			if speciesName != "" {
				counts[speciesName]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf(
			"An error occurred while reading the file: %v", err)
	}

	var sorted []KrakenSpecies
	for name, count := range counts {
		sorted = append(sorted, KrakenSpecies{Name: name, Count: count})
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Count == sorted[j].Count {
			return sorted[i].Name < sorted[j].Name
		}
		return sorted[i].Count > sorted[j].Count
	})

	var first, second *KrakenSpecies
	if len(sorted) > 0 {
		first = &sorted[0]
	}
	if len(sorted) > 1 {
		second = &sorted[1]
	}

	return first, second, nil
}
