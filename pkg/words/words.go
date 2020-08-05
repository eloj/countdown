package words

import (
	"io"
	"fmt"
	"os"
	"bufio"
	"sort"
	"math/bits"
)

// TODO: Add compare function for these, use instead of validateWord?
type wordentry struct {
	key  uint32
	dups int 		// number of duplicate characters in word.
	word string
	sorted string
}

func NewWordEntry(word string) wordentry {
	we := wordentry{}

	// TODO: Trim?
	we.word = NormalizeLatin1(word)
	we.sorted = sortWord(we.word)
	we.key = deriveKey32(we.sorted)
	we.dups = len(word) - bits.OnesCount32(we.key)

	return we
}

func sortWord(word string) string {
	ra := []rune(word)
	sort.Slice(ra, func(i, j int) bool { return ra[i] < ra[j] })
	return string(ra)
}

// TODO: Should have re-mapping support for extra chars. Room for 6 (32-25-1)
func deriveKey32(word string) uint32 {
	var key uint32

	for i:=0 ; i < len(word) ; i++ {
		ch := word[i]
		if (ch >= 'a' && ch <= 'z') {
			key |= 1 << (ch - 'a')
		} else {
			key |= 1 << 31 // Unknown-bit
		}
	}
	return key
}

type Countdown struct {
	minlen int
	maxlen int
	words []wordentry
}

func NewCountdown(minlen int, maxlen int) (*Countdown) {
	cd := &Countdown{}
	cd.minlen = minlen
	cd.maxlen = maxlen
	cd.words = make([]wordentry, 0, 1024)
	return cd
}


// All characters in word must exist in target.
func verifyWord(word string, target string) bool {
	sw := word
	tw := target

	var j int
	var i int
	for i = 0 ; i < len(sw) && j < len(tw) ; i++ {
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

func (cd *Countdown) FindWords(s string, maxdist int) {
	target := NewWordEntry(s)

	fmt.Printf("FIND on %#v\n", target)

	var hits int
	var numFalse int
	var numInvalid int
	var numDistFail int
	for i, word := range cd.words {
		if (len(word.word) < cd.minlen) || (len(word.word) > cd.maxlen) {
			continue
		}

		falsebits := bits.OnesCount32((target.key ^ word.key) & word.key)

		// fmt.Printf("I:%032b\nW:%032b -- %s\n= %032b (false=%d)\n", target.key, word.key, word.word, target.key & word.key, falsebits)

		if falsebits == 0 {
			hamming_est := len(target.sorted)- bits.OnesCount32(target.key & word.key)

			if /* hamming_est <= maxdist && */ verifyWord(word.sorted, target.sorted) {
				dist := len(target.sorted) - len(word.sorted)
				if (dist <= maxdist) {
					fmt.Printf("Found word #%d '%s', hamming estimate=%d, real distance=%d\n", i, word.word, hamming_est, dist)
					hits++
				} else {
					numDistFail++
				}
			} else {
				numInvalid++
			}
		} else {
			numFalse++
		}
	}
	fmt.Printf("%d words found, %d rejected by falsebits, %d rejected in validation, %d rejected by distance.\n", hits, numFalse, numInvalid, numDistFail)
}

func (cd *Countdown) addWord(word string) bool {
	we := NewWordEntry(word)

	if (we.key == 0) {
		return false
	}

	cd.words = append(cd.words, we)

	return true
}

func (cd *Countdown) AddDictionary(r io.Reader) (int, error) {
	var err error
	var cnt int

	src := bufio.NewScanner(r) // This has a line length limit, but that's okay for our use-case.

	for src.Scan() {
		line := src.Text()
		if (len(line) >= cd.minlen) && (len(line) <= cd.maxlen) && cd.addWord(line) {
			cnt++
		}
	}

	return cnt, err
}

func (cd *Countdown) AddDictionaryFile(filename string) (int, error) {
	var fh *os.File
	var cnt int

	fh, err := os.Open(filename)
	if err == nil {
		cnt, err = cd.AddDictionary(fh)
		fh.Close()
	}

	return cnt, err
}
