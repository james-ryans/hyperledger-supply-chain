package usermodel

import "time"

type ScanHistory struct {
	ID         string    `json:"id"`
	RiceSackID string    `json:"rice_sack_id"`
	ScanAt     time.Time `json:"scan_at"`
}
