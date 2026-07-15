package pipeline

import (
	"strings"
)

type SpeciesResult struct {
	DisplayName    string
	MLSTSpecies    string
	OtherMutations []string
	PoliMutations  []string
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func containsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func isEnterobacter(normalizedName string) bool {
	species := []string{
		"enterobactercloacae", "enterobacterasburiae",
		"enterobacterbugandensis", "enterobactercancerogenus",
		"enterobacterchengduensis", "enterobacterhormaechei",
		"enterobacterkobei", "enterobacterludwigii", "enterobactermori",
		"enterobacterroggenkampii", "enterobactersichuanensis",
		"enterobactersoli",
	}

	return containsAny(normalizedName, species)
}

func isAcinetobacter(normalizedName string) bool {
	species := []string{
		"acinetobacterbaumannii", "acinetobactercalcoaceticus",
		"acinetobacterlactucae", "acinetobacterpittii",
		"acinetobacterseifertii", "acinetobacternosocomialis",
	}

	return containsAny(normalizedName, species)
}

func isKlebsiella(normalizedName string) bool {
	species := []string{
		"klebsiellapneumoniae",
	}

	return containsAny(normalizedName, species)
}

func isPseudomonas(normalizedName string) bool {
	species := []string{
		"pseudomonasaeruginosa",
	}

	return containsAny(normalizedName, species)
}
