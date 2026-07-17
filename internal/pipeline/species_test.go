package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapitalizeFirst(t *testing.T) {
	t.Run("Lowercase", func(t *testing.T) {
		assert.Equal(t, "Escherichia", capitalizeFirst("escherichia"))
	})

	t.Run("Already Capitalized", func(t *testing.T) {
		assert.Equal(t, "Escherichia", capitalizeFirst("Escherichia"))
	})

	t.Run("Uppercase Input", func(t *testing.T) {
		assert.Equal(t, "Klebsiella", capitalizeFirst("KLEBSIELLA"))
	})

	t.Run("Single Character", func(t *testing.T) {
		assert.Equal(t, "E", capitalizeFirst("e"))
		assert.Equal(t, "E", capitalizeFirst("E"))
	})

	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, "", capitalizeFirst(""))
	})
}

func TestContainsAny(t *testing.T) {
	t.Run("First Substring", func(t *testing.T) {
		assert.True(t, containsAny("klebsiellapneumoniae", []string{"klebsiellapneumoniae"}))
	})

	t.Run("Later Substring", func(t *testing.T) {
		assert.True(t, containsAny("acinetobacterbaumannii", []string{"x", "y", "acinetobacterbaumannii"}))
	})

	t.Run("Partial Match", func(t *testing.T) {
		assert.True(t, containsAny("enterobactercloacae complex", []string{"enterobactercloacae"}))
	})

	t.Run("No Match", func(t *testing.T) {
		assert.False(t, containsAny("ecoli", []string{"klebsiellapneumoniae", "pseudomonasaeruginosa"}))
	})

	t.Run("Empty List", func(t *testing.T) {
		assert.False(t, containsAny("anything", []string{}))
	})

	t.Run("Empty Source", func(t *testing.T) {
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

	t.Run("E. kobei", func(t *testing.T) {
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
