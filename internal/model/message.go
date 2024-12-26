package model

type Message struct {
	ID      int    `json:"id"`
	Message string `json:"message,omitempty"`
}
