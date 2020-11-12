package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <listen-addr>", os.Args[0])
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	log.Fatal(http.ListenAndServe(os.Args[1], nil))
}
