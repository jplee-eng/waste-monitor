package main

import (
	"log"
	"net/http"
	"waste-monitor/internal/handler"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	h, err := handler.New()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", h.ServeHtml)
	http.HandleFunc("/styles.css", h.ServeCss)
	http.HandleFunc("/script.js", h.ServeJs)
	http.HandleFunc("/api/reading", h.HandleNewReading)
	http.HandleFunc("/events", h.HandleSSE)
	log.Println("up:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
