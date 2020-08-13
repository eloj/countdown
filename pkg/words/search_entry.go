package words

import (
	"sort"
	"strings"
)

// CountdownSearchEntry represents a single word.
type CountdownSearchEntry struct {
	key    uint32
	word   string
	sorted string
	// dups   int // number of duplicate characters in word.
}

// NewCountdownSearchEntry takes a word and generates a search-entry from it.
func NewCountdownSearchEntry(word string) CountdownSearchEntry {
	we := CountdownSearchEntry{}

	we.word = strings.TrimSpace(NormalizeLatin1(word))
	we.sorted = sortWord(we.word)
	we.key = deriveKey32(we.sorted)
	// we.dups = len(word) - bits.OnesCount32(we.key)

	return we
}

// Helper function to sort a string.
func sortWord(word string) string {
	ra := []rune(word)
	sort.Slice(ra, func(i, j int) bool { return ra[i] < ra[j] })
	return string(ra)
}

// Derive a key by setting bits based on which characters are in the word.
func deriveKey32(word string) uint32 {
	var key uint32

	for i := 0; i < len(word); i++ {
		ch := word[i]
		if ch >= 'a' && ch <= 'z' {
			key |= 1 << (ch - 'a')
		} else {
			key |= 1 << 31 // Unknown-bit
		}
	}
	return key
}

// Returns true iff all characters in word exist in target.
func verifyWord(word string, target string) bool {
	sw := word
	tw := target

	var j int
	var i int
	for i = 0; i < len(sw) && j < len(tw); i++ {
		// Scan forward if current character larger than target's.
		for sw[i] > tw[j] {
			j++
			if j >= len(tw) {
				break
			}
		}
		// Match, step target forward once.
		if sw[i] == tw[j] {
			j++
		} else {
			return false
		}
	}
	// Verify all input matched.
	return i == len(sw)
}
