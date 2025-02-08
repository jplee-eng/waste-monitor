package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
	"waste-monitor/internal/model"
)

func main() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		reading := model.Reading{
			Level:   rand.Intn(100),
			Battery: 3.0 + rand.Float64(),
			RSSI:    -110 + rand.Intn(30),
		}

		data, _ := json.Marshal(reading)
		resp, err := http.Post("http://localhost:8080/api/reading",
			"application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		resp.Body.Close()
		log.Printf("sent reading: %+v", reading)
	}
}
