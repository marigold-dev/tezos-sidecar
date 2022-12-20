package main

import (
	"encoding/json"
	"io/ioutil"
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
	tezosURI := env.GetDefault("TEZOS_URI", "https://mainnet.tezos.marigold.dev")

	// Health endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		level, err := request(tezosURI)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// If the new level is greater than the last level
		// or the timestamp is less than 5 minutes ago
		// return 200
		if level.Level > lastBlock.Level || level.Timestamp.Sub(time.Now()) <= (5*time.Minute) {
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var block BlockHeader
	json.Unmarshal(body, &block)

	log.Printf("Level: %d, Timestamp: %s", block.Level, block.Timestamp)
	result := &BlockSnapshot{Level: block.Level, Timestamp: time.Now()}
	return result, nil
}
