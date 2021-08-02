package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NotSoFancyName/conversion_service/service/fetch"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	stop := make(chan struct{})

	fetcher, err := fetch.NewFetcher(10 * time.Minute)
	if err != nil {
		log.Fatalf("failed to create a fetcher service: %v", err)
	}

	go fetcher.Run(stop)
	<-done
	log.Println("Request to stop the fetcher service")
	stop <- struct{}{}
	<-stop
}
