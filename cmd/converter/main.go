package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NotSoFancyName/conversion_service/service/rest"
)

var port = flag.Int("p", 8081, "listen port number")

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	errs := make(chan error)
	stop := make(chan struct{})

	srv, err := rest.NewServer(stop, *port)
	if err != nil {
		log.Fatalf("Failed to initialize converter server: %v", err)
	}
	go srv.Run(errs)
	select {
	case err := <-errs:
		log.Printf("Converter server stopped due to error: %v", err)
	case <-done:
		stop <- struct{}{}
		log.Println("Request to stop the coverter server")
		<-stop
	}
}
