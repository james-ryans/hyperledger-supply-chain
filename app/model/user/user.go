package usermodel

type User struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Email         string        `json:"email"`
	ScanHistories []ScanHistory `json:"scan_histories"`
}
