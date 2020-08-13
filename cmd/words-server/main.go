package main

/*
	Serves requests at /countdown/v1/words/<string>
*/

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	cntw "countdown/pkg/words"
)

type State struct {
	cw *cntw.Countdown
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}

type WordsResponse struct {
	Query      string   `json:"query"`
	Duration   float64  `json:"duration"` // In ms, document.
	NumHits    int      `json:"num_hits"`
	NumChecked int      `json:"num_checked"`
	MinDist    int      `json:"min_dist,omitempty"` // Only valid if words > 0
	MaxDist    int      `json:"max_dist,omitempty"` // ibid.
	Words      []string `json:"words"`
}

var (
	conf  Config
	state State
)

func SendResponse(w http.ResponseWriter, status int, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

func SendErrorResponse(w http.ResponseWriter, status int, msg string, desc string) error {
	res := ErrorResponse{Error: msg, Description: desc}
	return SendResponse(w, status, res)

}

func (state *State) wordsHandler(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.URL.Path, "/countdown/v1/words/") {
		w.WriteHeader(http.StatusNotFound)
		return

	}

	switch r.Method {
	case "GET":
		// TODO: Separate module for parsing and validation + error returns.
		// args, err := url.ParseQuery(r.URL.Query())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	scramble := strings.TrimPrefix(r.URL.Path, "/countdown/v1/words/")

	// TODO: Limit parameter ranges. Error on out-of-range.
	maxdist := 4
	maxhits := 10

	start := time.Now()
	result := state.cw.FindWords(scramble, maxhits, maxdist)
	elapsed := time.Since(start)

	// Sort and extract just the words for the response
	sorted := result.Sort()

	sorted_words := make([]string, len(sorted))
	for i, worddist := range sorted {
		sorted_words[i] = worddist.Word
	}

	res := WordsResponse{
		Query:      result.Query,
		Duration:   float64(elapsed) / float64(time.Millisecond),
		NumHits:    result.NumHits,
		NumChecked: result.NumChecked,
		MinDist:    0,
		MaxDist:    0,
		Words:      sorted_words,
	}

	// Extract actual min and max distance of result
	if len(sorted_words) > 0 {
		res.MinDist = sorted[0].Dist
		res.MaxDist = sorted[len(sorted_words)-1].Dist
	}

	SendResponse(w, 200, res)
}

func startHttpServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", conf.Port),
	}

	http.HandleFunc("/countdown/v1/words/", state.wordsHandler)

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// TODO: Log the error
			fmt.Fprintf(os.Stderr, "ListenAndServer: %v", err)
		}
	}()

	return srv
}

func main() {
	conffile := flag.String("config", "config/words-server.yaml", "Configuration file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if err := conf.ReadConfigurationFile(*conffile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read configuration from '%s': %v\n", *conffile, err)
		os.Exit(1)
	}

	fmt.Printf("CONFIG:%#v\n", conf)

	state.cw = cntw.NewCountdown(conf.MinWordLen, conf.MaxWordLen)

	for _, dict := range conf.Dictionaries {
		num, err := state.cw.AddDictionaryFile(dict)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading dictionary '%s', skipped.\n", dict)
		} else {
			fmt.Printf("%d words loaded from '%s'.\n", num, dict)
		}
	}

	var httpServerWg sync.WaitGroup
	srv := startHttpServer(&httpServerWg)

	signalch := make(chan os.Signal, 1)
	signal.Notify(signalch, os.Interrupt, syscall.SIGTERM)

	fmt.Printf("Server running. Waiting for signal.\n")

	select {
	case <-signalch:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		// log error
	}

	httpServerWg.Wait()

	fmt.Printf("All done.\n")
}
