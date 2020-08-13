package words

import (
	"bufio"
	"io"
	"math/bits"
	"os"
	"sort"
)

// clampInt Clamps an int to the range determined by the lo and hi arguments and returns it.
func clampInt(value int, lo int, hi int) int {
	if value < lo {
		value = lo
	} else if value > hi {
		value = hi
	}
	return value
}

// Countdown represents a dictionary of words of some minimum and maximum length,
// and partitioned by length.
type Countdown struct {
	minlen int
	maxlen int
	lvl    []CountdownWords
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
	levels := maxlen - minlen + 1
	cd.lvl = make([]CountdownWords, levels)
	for i := range cd.lvl {
		cd.lvl[i] = NewCountdownWords()
	}
	return cd
}

func (cd *Countdown) FindWords(s string, limit int, maxdist int) FindWordsResult {
	target := NewCountdownSearchEntry(s)

	// Adjust to be within valid range.
	if maxdist < 0 || maxdist > cd.maxlen {
		maxdist = cd.maxlen - cd.minlen
	}

	// Ensure that poorly constructed clients can't allocate too much memory.
	limit = clampInt(limit, 0, 1<<12)

	result := NewFindWordsResult(limit)
	result.Query = target.word

scankeys:
	// Scan words of decreasing length ...
	for level := 0; level <= maxdist; level++ {
		for idx, wordKey := range cd.lvl[level].keys {
			result.NumChecked++

			// We can immediately reject words using characters that are not in the target.
			falsebits := bits.OnesCount32((target.key ^ wordKey) & wordKey)

			// fmt.Printf("I:%032b\nW:%032b -- %s\n= %032b (false=%d)\n", target.key, word.key, word.word, target.key & word.key, falsebits)

			if falsebits == 0 {
				word := cd.lvl[level].words[idx]
				// TODO: Test out if hamming estimate is a useful (and correct) optimization.
				// hamming_est := len(target.sorted) - bits.OnesCount32(target.key&word.key) // then maxdist + math.Abs(target.dups - word.dups)

				if verifyWord(word.sorted, target.sorted) {
					dist := len(target.sorted) - len(word.sorted)
					if maxdist < 0 || dist <= maxdist {
						result.Words = append(result.Words, WordDistResult{word.word, dist})
						result.NumHits++
						if result.NumHits >= limit && limit > 0 {
							break scankeys
						}
					} else {
						// This triggers if the target is longer than maxlen, a sort of mismatch between dictionary and input.
						result.NumDistFail++
					}
				} else {
					result.NumInvalid++
				}
			} else {
				result.NumFalseBits++
			}
		}
	}

	return result
}

func (cd *Countdown) addWord(word string) bool {
	we := NewCountdownSearchEntry(word)

	if we.key == 0 {
		return false
	}

	// Add the search data to the hierarchy, based on its length.
	level := cd.maxlen - len(we.sorted)
	cd.lvl[level].Add(we)

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
