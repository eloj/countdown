package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	cntw "countdown/pkg/words"
)

func main() {
	var minlen int
	var maxlen int
	var maxdist int
	var limit int
	var scramble string

	dict := flag.String("file", "data/words-countdown.txt", "Dictionary to load")
	flag.IntVar(&minlen, "minlen", 4, "Minimum word length")
	flag.IntVar(&maxlen, "maxlen", 9, "Maximum word length")
	flag.IntVar(&maxdist, "maxdist", -1, "Maximum word distance for match (-1=none)")
	flag.IntVar(&limit, "limit", 0, "Maximum number of results to return (0=unlimited)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <input>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	scramble = "IBASELGSK"
	if flag.NArg() > 0 {
		scramble = flag.Arg(0)
	}

	cw := cntw.NewCountdown(minlen, maxlen)

	cnt, err := cw.AddDictionaryFile(*dict)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d words loaded from dictionary.\n", cnt)
	if cnt > 0 {
		fmt.Printf("Looking for solutions to '%s', minlen=%d, maxlen=%d, limit=%d, maxdist=%d\n", scramble, minlen, maxlen, limit, maxdist)
		start := time.Now()
		result := cw.FindWords(scramble, limit, maxdist)
		duration := time.Since(start)

		fmt.Printf("%d words found, %d rejected by falsebits, %d rejected in validation, %d rejected by distance. %d words checked. Duration: %v\n",
			result.NumHits, result.NumFalseBits, result.NumInvalid, result.NumDistFail, result.NumChecked, duration)

		for i, rec := range result.Sort() {
			fmt.Printf("%d. '%s' dist=%d\n", i, rec.Word, rec.Dist)
		}
	}
}
