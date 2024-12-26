package model

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Login     string `json:"login"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}
