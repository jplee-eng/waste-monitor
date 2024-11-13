package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"waste-monitor/internal/model"
)

type Handler struct {
	db        *sql.DB
	broadcast chan model.Reading
	clients   map[chan model.Reading]bool
	clientsMu sync.Mutex
}

func New() (*Handler, error) {
	db, err := sql.Open("sqlite3", "./waste_monitor.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS readings (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            level INTEGER NOT NULL,
            battery REAL NOT NULL,
            rssi INTEGER NOT NULL,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return nil, err
	}
	h := &Handler{
		db:        db,
		broadcast: make(chan model.Reading),
		clients:   make(map[chan model.Reading]bool),
	}
	go h.broadcaster()
	return h, nil
}

func (h *Handler) HandleNewReading(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var reading model.Reading
	if err := json.NewDecoder(r.Body).Decode(&reading); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := h.db.Exec(
		"INSERT INTO readings (level, battery, rssi) VALUES (?, ?, ?)",
		reading.Level, reading.Battery, reading.RSSI,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	reading.ID = id
	reading.Timestamp = time.Now()
	h.broadcast <- reading
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reading)
}

func (h *Handler) HandleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	clientChan := make(chan model.Reading)
	h.clientsMu.Lock()
	h.clients[clientChan] = true
	h.clientsMu.Unlock()
	defer func() {
		h.clientsMu.Lock()
		delete(h.clients, clientChan)
		h.clientsMu.Unlock()
	}()
	rows, err := h.db.Query("SELECT id, level, battery, rssi, timestamp FROM readings ORDER BY timestamp DESC LIMIT 1")
	if err == nil {
		defer rows.Close()
		if rows.Next() {
			var reading model.Reading
			rows.Scan(&reading.ID, &reading.Level, &reading.Battery, &reading.RSSI, &reading.Timestamp)
			data, _ := json.Marshal(reading)
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		}
	}
	for {
		select {
		case reading := <-clientChan:
			data, _ := json.Marshal(reading)
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (h *Handler) broadcaster() {
	for reading := range h.broadcast {
		h.clientsMu.Lock()
		for clientChan := range h.clients {
			select {
			case clientChan <- reading:
			default:
			}
		}
		h.clientsMu.Unlock()
	}
}

func (h *Handler) ServeHtml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/min/index.min.html")
}

func (h *Handler) ServeCss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "static/min/styles.min.css")
}

func (h *Handler) ServeJs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, "static/min/script.min.js")
}
