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

	cntw "github.com/eloj/countdown/pkg/words"

	"github.com/rs/zerolog/log"
)

type serverState struct {
	cw *cntw.Countdown
}

type errorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}

type wordsResponse struct {
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
	state serverState
)

func sendResponse(w http.ResponseWriter, status int, payload interface{}) error {
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

func sendErrorResponse(w http.ResponseWriter, status int, msg string, desc string) error {
	res := errorResponse{Error: msg, Description: desc}
	return sendResponse(w, status, res)

}

func (state *serverState) wordsHandler(w http.ResponseWriter, r *http.Request) {
	// Setup defaults
	args := &wordParams{
		Limit:   conf.DefaultLimit,
		Maxdist: conf.DefaultMaxDist,
	}

	if !strings.HasPrefix(r.URL.Path, "/countdown/v1/words/") {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		if err := parseWordParams(r.URL, args); err != nil {
			log.Error().Err(err).Msg("Query parameter error")
			sendErrorResponse(w, http.StatusBadRequest, "Query parameter error", err.Error())
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	scramble := strings.TrimPrefix(r.URL.Path, "/countdown/v1/words/")

	start := time.Now()
	result := state.cw.FindWords(scramble, args.Limit, args.Maxdist)
	elapsed := time.Since(start)

	// Sort and extract just the words for the response
	sorted := result.Sort()

	sortedWords := make([]string, len(sorted))
	for i, worddist := range sorted {
		sortedWords[i] = worddist.Word
	}

	durationMs := float64(elapsed) / float64(time.Millisecond)

	res := wordsResponse{
		Query:      result.Query,
		Duration:   durationMs,
		NumHits:    result.NumHits,
		NumChecked: result.NumChecked,
		MinDist:    0,
		MaxDist:    0,
		Words:      sortedWords,
	}

	// Extract actual min and max distance of result
	if len(sortedWords) > 0 {
		res.MinDist = sorted[0].Dist
		res.MaxDist = sorted[len(sortedWords)-1].Dist
	}

	log.Debug().Interface("args", args).Str("q", scramble).Int("hits", res.NumHits).Int("checked", res.NumChecked).Dur("find_ms", elapsed).Msg("Query")

	sendResponse(w, 200, res)
}

func startHTTPServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", conf.Port),
	}

	http.HandleFunc("/countdown/v1/words/", state.wordsHandler)

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error().Err(err).Msg("ListenAndServe error")
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

	if err := conf.readConfigurationFile(*conffile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read configuration from '%s': %v\n", *conffile, err)
		os.Exit(1)
	}

	log.Info().Interface("configuration", conf).Msg("")

	state.cw = cntw.NewCountdown(conf.MinWordLen, conf.MaxWordLen)

	for _, dict := range conf.Dictionaries {
		num, err := state.cw.AddDictionaryFile(dict)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading dictionary '%s', skipped.\n", dict)
		} else {
			log.Info().Int("words", num).Str("source", dict).Msg("Loaded dictionary.")
		}
	}

	var httpServerWg sync.WaitGroup
	srv := startHTTPServer(&httpServerWg)

	signalch := make(chan os.Signal, 1)
	signal.Notify(signalch, os.Interrupt, syscall.SIGTERM)

	log.Info().Msg("Server started. Waiting for signal.")

	select {
	case <-signalch:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Shutdown error")
	}

	httpServerWg.Wait()

	log.Info().Msg("Server exiting.")
}
