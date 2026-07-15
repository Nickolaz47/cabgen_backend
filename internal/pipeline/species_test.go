package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapitalizeFirst(t *testing.T) {
	t.Run("Capitalizes lowercase word", func(t *testing.T) {
		assert.Equal(t, "Escherichia", capitalizeFirst("escherichia"))
	})

	t.Run("Capitalizes already capitalized", func(t *testing.T) {
		assert.Equal(t, "Escherichia", capitalizeFirst("Escherichia"))
	})

	t.Run("Lowercases the rest", func(t *testing.T) {
		assert.Equal(t, "Klebsiella", capitalizeFirst("KLEBSIELLA"))
	})

	t.Run("Single character", func(t *testing.T) {
		assert.Equal(t, "E", capitalizeFirst("e"))
		assert.Equal(t, "E", capitalizeFirst("E"))
	})

	t.Run("Empty string", func(t *testing.T) {
		assert.Equal(t, "", capitalizeFirst(""))
	})
}

func TestContainsAny(t *testing.T) {
	t.Run("Match found in first substring", func(t *testing.T) {
		assert.True(t, containsAny("klebsiellapneumoniae", []string{"klebsiellapneumoniae"}))
	})

	t.Run("Match found in later substring", func(t *testing.T) {
		assert.True(t, containsAny("acinetobacterbaumannii", []string{"x", "y", "acinetobacterbaumannii"}))
	})

	t.Run("Partial match via Contains", func(t *testing.T) {
		assert.True(t, containsAny("enterobactercloacae complex", []string{"enterobactercloacae"}))
	})

	t.Run("No match", func(t *testing.T) {
		assert.False(t, containsAny("ecoli", []string{"klebsiellapneumoniae", "pseudomonasaeruginosa"}))
	})

	t.Run("Empty substrings list", func(t *testing.T) {
		assert.False(t, containsAny("anything", []string{}))
	})

	t.Run("Empty source string", func(t *testing.T) {
		assert.False(t, containsAny("", []string{"x"}))
	})
}

func TestIsEnterobacter(t *testing.T) {
	t.Run("E. cloacae", func(t *testing.T) {
		assert.True(t, isEnterobacter("enterobactercloacae"))
	})

	t.Run("E. hormaechei", func(t *testing.T) {
		assert.True(t, isEnterobacter("enterobacterhormaechei"))
	})

	t.Run("E. kobei as substring", func(t *testing.T) {
		assert.True(t, isEnterobacter("enterobacterkobei strain"))
	})

	t.Run("Non-enterobacter", func(t *testing.T) {
		assert.False(t, isEnterobacter("escherichiacoli"))
	})
}

func TestIsAcinetobacter(t *testing.T) {
	t.Run("A. baumannii", func(t *testing.T) {
		assert.True(t, isAcinetobacter("acinetobacterbaumannii"))
	})

	t.Run("A. pittii", func(t *testing.T) {
		assert.True(t, isAcinetobacter("acinetobacterpittii"))
	})

	t.Run("A. nosocomialis", func(t *testing.T) {
		assert.True(t, isAcinetobacter("acinetobacternosocomialis"))
	})

	t.Run("Non-acinetobacter", func(t *testing.T) {
		assert.False(t, isAcinetobacter("pseudomonasaeruginosa"))
	})
}

func TestIsKlebsiella(t *testing.T) {
	t.Run("K. pneumoniae", func(t *testing.T) {
		assert.True(t, isKlebsiella("klebsiellapneumoniae"))
	})

	t.Run("K. pneumoniae with extra text", func(t *testing.T) {
		assert.True(t, isKlebsiella("klebsiellapneumoniae subsp"))
	})

	t.Run("Non-klebsiella", func(t *testing.T) {
		assert.False(t, isKlebsiella("enterobactercloacae"))
	})
}

func TestIsPseudomonas(t *testing.T) {
	t.Run("P. aeruginosa", func(t *testing.T) {
		assert.True(t, isPseudomonas("pseudomonasaeruginosa"))
	})

	t.Run("P. aeruginosa with extra text", func(t *testing.T) {
		assert.True(t, isPseudomonas("pseudomonasaeruginosa strain"))
	})

	t.Run("Non-pseudomonas", func(t *testing.T) {
		assert.False(t, isPseudomonas("acinetobacterbaumannii"))
	})
}