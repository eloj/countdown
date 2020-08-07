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

var (
	conf  Config
	state State
)

func Response(w http.ResponseWriter, status int, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

func (state *State) wordsHandler(w http.ResponseWriter, r *http.Request) {
	scramble := strings.TrimPrefix(r.URL.Path, "/countdown/v1/words/")

	maxdist := 3

	result := state.cw.FindWords(scramble, maxdist)

	Response(w, 200, result.Sort())
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
