package words

// CountdownWords represents a set search words and keys sharing the same length.
type CountdownWords struct {
	keys  []uint32 // The keys array more than doubles the search speed vs iterating over the words array directly.
	words []CountdownSearchEntry
}

func NewCountdownWords() CountdownWords {
	cdw := CountdownWords{}
	cdw.keys = make([]uint32, 0, 1024)
	cdw.words = make([]CountdownSearchEntry, 0, 1024)
	return cdw
}

func (cdw *CountdownWords) Add(se CountdownSearchEntry) {
	cdw.keys = append(cdw.keys, se.key)
	cdw.words = append(cdw.words, se)
}
