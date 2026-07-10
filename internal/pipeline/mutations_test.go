package pipeline

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createMockBlastFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "blast_mock_*.txt")
	assert.NoError(t, err)

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}

const organismMockContent = `
> GyrA|
Length=100
Identities = 95/95 (95%)
Query  1   ATCG 4
           || |
Sbjct  1   ATAG 4

> PmrA|
Length=100
Identities = 95/95 (95%)
Query  1   ATCG 4
           || |
Sbjct  1   ATAG 4
`

func TestFindMutation(t *testing.T) {
	mockContent := `
> GyrA|
Length=100
Identities = 95/95 (95%)
Query  1   ATCG 4
           || |
Sbjct  1   ATAG 4

> PmrA|
Length=100
Identities = 85/85 (85%)
`
	path := createMockBlastFile(t, mockContent)

	finder := NewMutationFinder(path).(*mutationFinder)

	t.Run("Success - Extraction and Truncation", func(t *testing.T) {
		res, err := finder.findMutation([]string{"GyrA", "PmrA"})
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Contains(t, res, "GyrA:A3C,")
		assert.Contains(t, res, "PmrA truncation: 85/100,")
	})

	t.Run("Error - Empty File", func(t *testing.T) {
		emptyPath := createMockBlastFile(t, "")

		emptyFinder := NewMutationFinder(emptyPath).(*mutationFinder)
		res, err := emptyFinder.findMutation([]string{"GyrA"})
		assert.Error(t, err)
		assert.Equal(t, "Empty BLAST result", err.Error())
		assert.Nil(t, res)
	})

	t.Run("Error - File Not Found", func(t *testing.T) {
		notFoundFinder := NewMutationFinder("invalid_path_to_blast.txt").(*mutationFinder)
		res, err := notFoundFinder.findMutation([]string{"GyrA"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "BLAST result file not found")
		assert.Nil(t, res)
	})
}

func TestAcinetoMutations(t *testing.T) {
	path := createMockBlastFile(t, organismMockContent)

	finder := NewMutationFinder(path)
	other, poli, err := finder.FindAcinetoMutations()

	assert.NoError(t, err)
	assert.Contains(t, other, "GyrA:A3C,")
	assert.Contains(t, poli, "PmrA:A3C,")
}

func TestEcloacaeMutations(t *testing.T) {
	path := createMockBlastFile(t, organismMockContent)

	finder := NewMutationFinder(path)
	other, poli, err := finder.FindEcloacaeMutations()

	assert.NoError(t, err)
	assert.Contains(t, other, "GyrA:A3C,")
	assert.Contains(t, poli, "PmrA:A3C,")
}

func TestKlebMutations(t *testing.T) {
	path := createMockBlastFile(t, organismMockContent)

	finder := NewMutationFinder(path)
	other, poli, err := finder.FindKlebMutations()

	assert.NoError(t, err)
	assert.Contains(t, other, "GyrA:A3C,")
	assert.Contains(t, poli, "PmrA:A3C,")
}

func TestPseudoMutations(t *testing.T) {
	path := createMockBlastFile(t, organismMockContent)

	finder := NewMutationFinder(path)
	other, poli, err := finder.FindPseudoMutations()

	assert.NoError(t, err)
	assert.Contains(t, other, "GyrA:A3C,")
	assert.Contains(t, poli, "PmrA:A3C,")
}
