package model

import "time"

type Reading struct {
	ID        int64     `json:"id"`
	Level     int       `json:"level"`
	Battery   float64   `json:"battery"`
	RSSI      int       `json:"rssi"`
	Timestamp time.Time `json:"timestamp"`
}
