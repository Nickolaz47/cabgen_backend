package pipeline

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type KrakenSpecies struct {
	Name  string
	Count int
}

func KrakenSpeciesCounter(krakenReport string) (*KrakenSpecies, *KrakenSpecies,
	error) {
	file, err := os.Open(krakenReport)
	if err != nil {
		return nil, nil, fmt.Errorf("Kraken report file not found: %v", err)
	}
	defer file.Close()

	br := bufio.NewReader(file)
	_, err = br.Peek(1)
	if err == io.EOF {
		return nil, nil, errors.New("Empty Kraken report")
	}

	scanner := bufio.NewScanner(br)
	const maxCapacity = 1024 * 1024 * 50
	lineBuf := make([]byte, 1024*64)
	scanner.Buffer(lineBuf, maxCapacity)

	var species []KrakenSpecies

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) < 6 {
			continue
		}

		if strings.TrimSpace(fields[3]) != "S" {
			continue
		}

		cladeReads, err := strconv.Atoi(strings.TrimSpace(fields[1]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(fields[5])
		if name == "" {
			continue
		}

		species = append(species, KrakenSpecies{Name: name, Count: cladeReads})
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading kraken2 report: %w", err)
	}

	sort.Slice(species, func(i, j int) bool {
		if species[i].Count == species[j].Count {
			return species[i].Name < species[j].Name
		}
		return species[i].Count > species[j].Count
	})

	var first, second *KrakenSpecies
	if len(species) > 0 {
		first = &species[0]
	}
	if len(species) > 1 {
		second = &species[1]
	}

	return first, second, nil
}
