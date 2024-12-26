package model

type SubscriberList struct {
	Subscribers []Subscriber `json:"subscribers,omitempty"`
	TotalCount  int          `json:"total_count"`
	ActiveCount int          `json:"active_count,omitempty"`
}
