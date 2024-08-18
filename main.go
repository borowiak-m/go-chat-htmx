package main

import (
	"log"
	"net/http"

	"github.com/borowiak-m/go-chat-htmx/types"
)

func serveIndex(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if req.Method != "GET" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, req, "templates/index.html")
}

func main() {
	hub := types.NewHub()
	go hub.Run()
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		types.ServeWs(hub, w, req)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
