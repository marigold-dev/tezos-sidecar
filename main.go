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
	tezosURI := env.GetDefault("TEZOS_URI", "https://kathmandunet.tezos.marigold.dev")

	// Health endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		level, err := request(tezosURI)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if level.Level < lastBlock.Level && level.Timestamp.Sub(time.Now()) < (5*time.Minute) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		lastBlock = *level
		w.WriteHeader(http.StatusOK)
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
