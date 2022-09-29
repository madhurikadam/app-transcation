package domain

import "time"

type Account struct {
	ID             string     `json:"id"`
	DocumentNumber string     `json:"document_number"`
	Balance        float64    `json:"balance"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type AccountReq struct {
	DocumentNumber string `json:"document_number"`
}

type Transcation struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventAt         time.Time `json:"event_at"`
}
