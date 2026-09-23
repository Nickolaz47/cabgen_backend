package validations_test

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var exemptedFromSanitize = map[string]bool{
	"AdminAnalysisCreateInput": true,
	"AnalysisCreateInput":      true,
	"AnalysisTSVDownloadInput": true,
	"UpdatePasswordInput":      true,
}

func readValidationsSource(t *testing.T) string {
	src, err := os.ReadFile("../validations/validations.go")
	assert.NoError(t, err)
	return string(src)
}

func extractModelUnionTypes(t *testing.T, src string) []string {
	start := strings.Index(src, "type Model interface {")
	assert.NotEqual(t, -1, start, "Model interface not found")
	end := strings.Index(src[start:], "}")
	assert.NotEqual(t, -1, end, "Model interface body not found")

	block := src[start : start+end]
	re := regexp.MustCompile(`models\.(\w+)`)
	matches := re.FindAllStringSubmatch(block, -1)

	types := make([]string, 0, len(matches))
	for _, match := range matches {
		types = append(types, match[1])
	}
	sort.Strings(types)
	return types
}

func extractSanitizeSwitchCases(t *testing.T, src string) map[string]bool {
	switchStart := strings.Index(src, "func SanitizeInput")
	assert.NotEqual(t, -1, switchStart)

	re := regexp.MustCompile(`case \*models\.(\w+)`)
	matches := re.FindAllStringSubmatch(src[switchStart:], -1)

	cases := make(map[string]bool, len(matches))
	for _, match := range matches {
		cases[match[1]] = true
	}
	return cases
}

func TestAllModelInputsAreSanitizedOrExempted(t *testing.T) {
	src := readValidationsSource(t)
	union := extractModelUnionTypes(t, src)
	cases := extractSanitizeSwitchCases(t, src)

	assert.NotEmpty(t, union, "Model union should not be empty")
	assert.NotEmpty(t, cases, "SanitizeInput switch should not be empty")

	var missing []string
	for _, inputType := range union {
		if cases[inputType] || exemptedFromSanitize[inputType] {
			continue
		}
		missing = append(missing, inputType)
	}

	assert.Empty(t, missing,
		"input types in the Model union must have a SanitizeInput case "+
			"or be in the exempted list: %s", strings.Join(missing, ", "))
}
