package usermodel

import "time"

type ScanHistory struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	RiceSackCode string    `json:"rice_sack_code"`
	Name         string    `json:"name"`
	ScanAt       time.Time `json:"scan_at"`
}

type ScanHistoryService interface {
	GetAllScanHistories(userID string) ([]*ScanHistory, error)
}

type ScanHistoryRepository interface {
	FindAll(userID string) ([]*ScanHistory, error)
}
