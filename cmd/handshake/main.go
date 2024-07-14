package main

import (
	"log"
	"time"

	"github.com/reonardoleis/btc-handshake/pkg/handshake"
)

var (
	address           = ":8333"
	maxMockWorkerRuns = 20
)

func mockWorker(maxRuns int) {
	for i := 0; i < maxRuns; i++ {
		log.Println("Running mock worker on iteration", i)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	go handshake.Handshake(address)

	// just to simulate that the program could run other tasks
	// while doing the handshake
	// if it needed the response of the handshake, a channel
	// could be used in order to communicate the handshake completion :)
	mockWorker(maxMockWorkerRuns)
}
