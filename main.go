package main

import (
	"flag"
	"log"

	"github.com/RajNykDhulapkar/go-ww-sse/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":8080", "listen address")
	flag.Parse()

	server := api.NewServer(*listenAddr)
	log.Println("Starting server on", *listenAddr)
	log.Fatal(server.Start())
}
