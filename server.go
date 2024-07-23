package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":1234", "http service address")

type Server struct {
	matches map[string]*Match
}

func NewServer() *Server {
	return &Server{
		matches: make(map[string]*Match),
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	log.Println("Not ws")
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	log.Println("Starting server...")
	flag.Parse()
	hub := NewHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hitting")
		serveWs(hub, w, r)
	})
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
