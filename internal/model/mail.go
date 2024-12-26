package model

import "time"

type Mail struct {
	ID          int       `json:"id"`
	To          []string  `json:"to"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	ContentType string    `json:"content_type,omitempty"`
	SentAt      time.Time `json:"sent_at,omitempty"`
}
