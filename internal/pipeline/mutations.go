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

var (
	subjectPattern       = regexp.MustCompile(`(?i)>\s{0,1}(\w+)\|`)
	subjectLengthPattern = regexp.MustCompile(`(?i)Length=(\d+)`)
	identitiesPattern    = regexp.MustCompile(
		`(?i)Identities = \d+\/(\d+)\s\((\d+)%\)`)
	querySequencePattern = regexp.MustCompile(
		`(?i)Query\s+\d+\s+(\w+)`)
	subjectSequencePattern = regexp.MustCompile(
		`(?i)Sbjct\s+(\d+)\s+(\w+)\s+(\d+)`)
)

type MutationFinder interface {
	FindAcinetoMutations() ([]string, []string, error)
	FindEcloacaeMutations() ([]string, []string, error)
	FindKlebMutations() ([]string, []string, error)
	FindPseudoMutations() ([]string, []string, error)
}

type mutationFinder struct {
	BlastResultPath string
}

func NewMutationFinder(blastResultPath string) MutationFinder {
	return &mutationFinder{
		BlastResultPath: blastResultPath,
	}
}

func (f *mutationFinder) findMutation(mutations []string) (
	[]string, error) {
	file, err := os.Open(f.BlastResultPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read BLAST result: %v", err)
	}
	defer file.Close()

	br := bufio.NewReader(file)
	_, err = br.Peek(1)
	if err == io.EOF {
		return nil, errors.New("Empty BLAST result")
	}

	mutationsMap := make(map[string]bool)
	for _, t := range mutations {
		mutationsMap[t] = true
	}

	var foundMutations []string

	var subject string
	var subjectLength, identitiesTotal, percentualIdentity int
	var querySequence, subjectSequence, alignmentSequence []rune
	var subjectSequenceStart, subjectSequenceEnd int
	var test, nextLineIsAlignment bool

	scanner := bufio.NewScanner(br)

	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if matches := subjectPattern.FindStringSubmatch(line); len(
			matches) > 1 {
			subject = matches[1]
			test = false
			continue
		}

		if matches := subjectLengthPattern.FindStringSubmatch(line); len(
			matches) > 1 {
			subjectLength, _ = strconv.Atoi(matches[1])
			continue
		}

		if matches := identitiesPattern.FindStringSubmatch(line); len(
			matches) > 2 {
			identitiesTotal, _ = strconv.Atoi(matches[1])
			percentualIdentity, _ = strconv.Atoi(matches[2])

			if float64(identitiesTotal) < (float64(subjectLength)/100.0)*90.0 &&
				percentualIdentity > 80 {
				mutation := fmt.Sprintf(
					"%s truncation: %d/%d,", subject, identitiesTotal,
					subjectLength)
				foundMutations = append(foundMutations, mutation)
			}
			continue
		}

		if matches := querySequencePattern.FindStringSubmatch(line); len(
			matches) > 1 {
			if percentualIdentity > 90 {
				querySequence = []rune(strings.TrimSpace(matches[1]))
				nextLineIsAlignment = true
			}
			continue
		}

		if nextLineIsAlignment {
			alignmentSequence = []rune(strings.TrimSpace(line))
			nextLineIsAlignment = false
			continue
		}

		if matches := subjectSequencePattern.FindStringSubmatch(
			line); len(matches) > 3 {
			if percentualIdentity > 90 {
				subjectSequenceStart, _ = strconv.Atoi(matches[1])
				subjectSequence = []rune(strings.TrimSpace(matches[2]))
				subjectSequenceEnd, _ = strconv.Atoi(matches[3])

				if subjectSequenceStart == 1 {
					test = true
				}

				if test {
					for i := range 60 {
						if i >= len(querySequence) {
							break
						}

						if i >= len(alignmentSequence) || i >= len(
							subjectSequence) {
							break
						}

						queryAA := strings.ToUpper(string(querySequence[i]))
						alignmentAA := strings.ToUpper(string(alignmentSequence[i]))
						subjectAA := strings.ToUpper(string(subjectSequence[i]))

						if queryAA != alignmentAA {
							var position int
							if subjectSequenceStart > subjectSequenceEnd {
								position = subjectSequenceStart - i
							} else {
								position = subjectSequenceStart + i
							}

							if mutationsMap[subject] && queryAA != subjectAA {
								if float64(identitiesTotal) > (float64(
									subjectLength)/100.0)*90.0 {
									mutation := fmt.Sprintf(
										"%s:%s%d%s,", subject, subjectAA,
										position, queryAA)
									foundMutations = append(foundMutations,
										mutation)
								}
							}
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf(
			"An error occurred while reading the file: %v", err)
	}

	return foundMutations, nil
}

func (f *mutationFinder) FindAcinetoMutations() (
	[]string, []string, error) {
	otherMutations := []string{"GyrA", "GyrB", "ParC", "AdeN",
		"AdeR", "CarO", "OmpA", "AdeL", "AdeS"}
	poliMutations := []string{"PmrA", "PmrB", "LpxA", "LpxD", "LpxC"}

	otherResult, err := f.findMutation(otherMutations)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Failed to find other mutations to acineto: %v", err)
	}

	poliResult, err := f.findMutation(poliMutations)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Failed to find poli mutations to acineto: %v", err)
	}

	return otherResult, poliResult, nil
}

func (f *mutationFinder) FindEcloacaeMutations() (
	[]string, []string, error) {
	otherMutations := []string{"GyrA", "ParC"}
	poliMutations := []string{"PmrA", "PmrB", "MgrB", "PhoP", "PhoQ"}

	otherResult, err := f.findMutation(otherMutations)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Failed to find other mutations to ecloacae: %v", err)
	}

	poliResult, err := f.findMutation(poliMutations)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Failed to find poli mutations to ecloacae: %v", err)
	}

	return otherResult, poliResult, nil
}

func (f *mutationFinder) FindKlebMutations() (
	[]string, []string, error) {
	otherMutations := []string{"GyrA", "GyrB", "ParC", "AcrR", "RamR"}
	poliMutations := []string{"PmrB", "PmrA", "MgrB", "PhoP", "PhoQ"}

	otherResult, err := f.findMutation(otherMutations)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Failed to find other mutations to kleb: %v", err)
	}

	poliResult, err := f.findMutation(poliMutations)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Failed to find poli mutations to kleb: %v", err)
	}

	return otherResult, poliResult, nil
}

func (f *mutationFinder) FindPseudoMutations() (
	[]string, []string, error) {
	otherMutations := []string{"OprD", "MexT", "AmpC",
		"AmpR", "GyrA", "GyrB", "ParC", "ParE"}
	poliMutations := []string{"PmrA", "PmrB", "PhoQ",
		"ParR", "ParS", "CrpS", "ColR", "ColS"}

	otherResult, err := f.findMutation(otherMutations)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Failed to find other mutations to pseudo: %v", err)
	}

	poliResult, err := f.findMutation(poliMutations)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Failed to find poli mutations to pseudo: %v", err)
	}

	return otherResult, poliResult, nil
}
