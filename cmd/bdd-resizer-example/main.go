package main

import (
	"flag"
	"log"

	"github.com/hedhyw/bdd-resizer-example/internal/server"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9777", "address of the server")
	flag.Parse()

	s := server.New()

	log.Print("listening on ", *addr)

	err := s.ListenAndServer(*addr)
	if err != nil {
		log.Fatal(err)
	}
}
