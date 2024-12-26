package model

import "time"

// TODO add levels of subscribers
type Subscriber struct {
	ID                  int       `json:"id"`
	UserID              int       `json:"user_id"`
	StatusSubscription  string    `json:"status,omitempty"`
	NumberSubscriptions int       `json:"number_subscriptions"`
	SubscriptionTime    time.Time `json:"subscription_time"`
	SubscriptionsInRow  int       `json:"subscriptions_in_row"`
}
