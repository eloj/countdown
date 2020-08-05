package main

import (
	"fmt"
	"os"
	"log"
	"flag"

	cntw "countdown/pkg/words"
)

func main() {
	var minlen int
	var maxlen int
	var scramble string

	dict := flag.String("file", "data/words-countdown.txt", "Dictionary to load")
	flag.IntVar(&minlen, "minlen", 4, "Minimum word length")
	flag.IntVar(&maxlen, "maxlen", 9, "Maximum word length")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// dict := "/home/eddy/dev/english-words-alpha.txt"

	scramble = "SUDRELSIT"
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
		fmt.Printf("Looking for solutions to '%s', minlen=%d, maxlen=%d\n", scramble, minlen, maxlen)
		cw.FindWords(scramble, minlen)
	}
}
