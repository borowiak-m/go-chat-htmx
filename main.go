package main

import (
	"log"
	"net/http"
)

func serveIndex(w http.ResponseWriter, req *http.Request) {}

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", nil)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
