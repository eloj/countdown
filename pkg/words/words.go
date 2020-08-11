package words

import (
	"bufio"
	"io"
	"math/bits"
	"os"
	"sort"
)

// TODO: Add compare function for these, use instead of validateWord?
type searchEntry struct {
	key    uint32
	dups   int // number of duplicate characters in word.
	word   string
	sorted string
}

func NewWordEntry(word string) searchEntry {
	we := searchEntry{}

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

type Countdown struct {
	minlen int
	maxlen int
	words  []searchEntry
}

type WordDistResult struct {
	Word string
	Dist int
}

type FindWordsResult struct {
	Query string
	Words []WordDistResult

	NumChecked   int
	NumHits      int
	NumFalseBits int
	NumInvalid   int
	NumDistFail  int
}

func NewFindWordsResult(capacity int) FindWordsResult {
	result := FindWordsResult{}
	/*
		if capacity == 0 {
			capacity = 32
		}*/
	result.Words = make([]WordDistResult, 0, capacity)
	return result
}

func (result *FindWordsResult) Sort() []WordDistResult {
	sort.Slice(result.Words, func(i, j int) bool {
		if result.Words[i].Dist == result.Words[j].Dist {
			return result.Words[i].Word < result.Words[j].Word
		}
		return result.Words[i].Dist < result.Words[j].Dist
	})
	return result.Words
}

func NewCountdown(minlen int, maxlen int) *Countdown {
	cd := &Countdown{}
	cd.minlen = minlen
	cd.maxlen = maxlen
	cd.words = make([]searchEntry, 0, 1024)
	return cd
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

// ClampInt Clamps an int to the range determined by the lo and hi arguments and returns it.
// Can we have a math.Clamp() please?
func ClampInt(value int, lo int, hi int) int {
	if value < lo {
		value = lo
	} else if value > hi {
		value = hi
	}
	return value
}

func (cd *Countdown) FindWords(s string, maxhits int, maxdist int) FindWordsResult {
	target := NewWordEntry(s)

	// LOG: fmt.Printf("FIND on %#v\n", target)

	// Ensure that poorly constructed clients can't allocate too much memory.
	maxhits = ClampInt(maxhits, 0, 1<<12)

	result := NewFindWordsResult(maxhits)
	result.Query = target.word

	for _, word := range cd.words {
		result.NumChecked++

		falsebits := bits.OnesCount32((target.key ^ word.key) & word.key)

		// fmt.Printf("I:%032b\nW:%032b -- %s\n= %032b (false=%d)\n", target.key, word.key, word.word, target.key & word.key, falsebits)

		if falsebits == 0 {
			// TODO: Test out if hamming estimate is a useful (and correct) optimization.
			// hamming_est := len(target.sorted) - bits.OnesCount32(target.key&word.key) // then maxdist + math.Abs(target.dups - word.dups)

			if /* hamming_est <= maxdist && */ verifyWord(word.sorted, target.sorted) {
				dist := len(target.sorted) - len(word.sorted)
				if maxdist < 0 || dist <= maxdist {
					// fmt.Printf("Found word #%d '%s', hamming weight estimate=%d, real distance=%d\n", i, word.word, hamming_est, dist)
					result.Words = append(result.Words, WordDistResult{word.word, dist})
					result.NumHits++
				} else {
					result.NumDistFail++
				}
			} else {
				result.NumInvalid++
			}
		} else {
			result.NumFalseBits++
		}
	}

	return result
}

func (cd *Countdown) addWord(word string) bool {
	we := NewWordEntry(word)

	if we.key == 0 {
		return false
	}

	cd.words = append(cd.words, we)

	return true
}

func (cd *Countdown) AddDictionary(r io.Reader) (int, error) {
	var err error
	var cnt int

	src := bufio.NewScanner(r) // This has a line length limit, but that's okay for our use-case (else use .ReadString())

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
