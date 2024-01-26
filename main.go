package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	env "github.com/caitlinelfring/go-env-default"
)

type BlockSnapshot struct {
	Level     int
	Timestamp time.Time
}

type BlockHeader struct {
	Level     int       `json:"level"`
	Timestamp time.Time `json:"timestamp"`
}

var lastBlock BlockSnapshot = BlockSnapshot{Level: 0, Timestamp: time.Now()}

func main() {
	addr := env.GetDefault("ADDR", ":31234")
	minutes := env.GetIntDefault("MINUTES", 5)
	tezosURI := env.GetDefault("TEZOS_URI", "https://mainnet.tezos.marigold.dev")
	log.Printf("Listening on %s\n", addr)

	// Health endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		level, err := request(tezosURI)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// If the new level is greater than the last level
		// or the timestamp is less than X minutes ago
		// return 200
		if level.Level > lastBlock.Level || level.Timestamp.Sub(time.Now()) <= (time.Duration(minutes)*time.Minute) {
			lastBlock = *level
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func request(tezosURI string) (*BlockSnapshot, error) {
	resp, err := http.Get(tezosURI + "/chains/main/blocks/head/header")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var block BlockHeader
	json.Unmarshal(body, &block)

	log.Printf("Level: %d, Timestamp: %s", block.Level, block.Timestamp)
	result := &BlockSnapshot{Level: block.Level, Timestamp: block.Timestamp}
	return result, nil
}
