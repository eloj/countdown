package words

import (
	"sort"
	"strings"
)

type CountdownSearchEntry struct {
	key    uint32
	word   string
	sorted string
	// dups   int // number of duplicate characters in word.
}

func NewSearchEntry(word string) CountdownSearchEntry {
	we := CountdownSearchEntry{}

	we.word = strings.TrimSpace(NormalizeLatin1(word))
	we.sorted = sortWord(we.word)
	we.key = deriveKey32(we.sorted)
	// we.dups = len(word) - bits.OnesCount32(we.key)

	return we
}

func sortWord(word string) string {
	ra := []rune(word)
	sort.Slice(ra, func(i, j int) bool { return ra[i] < ra[j] })
	return string(ra)
}

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

// All characters in word must exist in target.
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
