package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NotSoFancyName/conversion_service/service/fetch"
)

const fetchInterval = 10 * time.Minute

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	stop := make(chan struct{})
	errs := make(chan error)

	fetcher, err := fetch.NewFetcher(fetchInterval)
	if err != nil {
		log.Fatalf("failed to create a fetcher service: %v", err)
	}

	go fetcher.RunFetcher(stop)
	go fetcher.RunRPCServer(":8082", errs)

	select {
	case <-done:
		log.Println("Request to stop the fetcher service")
	case <-errs:
		log.Println("Fatal error occured")
	}

	stop <- struct{}{}
	<-stop
}
