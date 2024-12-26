package subscriber

import (
	"context"
	"database/sql"
	"errors"
	"subscription-mailing-service/internal/config"
	"subscription-mailing-service/internal/model"
	"subscription-mailing-service/storage/postgres"
)

type SubscriberStorage struct {
	db *sql.DB
}

func (s *SubscriberStorage) Close() error {
	return postgres.CloseConnection(s.db)
}

func NewStorage(cfg *config.Config) (*SubscriberStorage, error) {
	db, err := postgres.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &SubscriberStorage{db: db}, nil
}

func (s *SubscriberStorage) Get(ctx context.Context, id int) (*model.Subscriber, error) {
	const query = `
		SELECT
		    user_id, 
		    status_subscription,
		    number_subscriptions, 
		    subscription_time,
		    subscriptions_in_row 
		FROM 
		    subscribers 
		WHERE 
		    id = $1
		    `

	subscriber := &model.Subscriber{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&subscriber.UserID,
		&subscriber.StatusSubscription,
		&subscriber.NumberSubscriptions,
		&subscriber.SubscriptionTime,
		&subscriber.SubscriptionsInRow,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return subscriber, err
}

func (s *SubscriberStorage) GetAll(ctx context.Context) ([]*model.Subscriber, error) {
	const query = `
		SELECT
		    id,
		    user_id, 
		    status_subscription,
		    number_subscriptions, 
		    subscription_time, 
		    subscriptions_in_row
		FROM
		    subscribers
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribers []*model.Subscriber
	for rows.Next() {
		var subscriber model.Subscriber
		if err := rows.Scan(
			&subscriber.ID,
			&subscriber.UserID,
			&subscriber.StatusSubscription,
			&subscriber.NumberSubscriptions,
			&subscriber.SubscriptionTime,
			&subscriber.SubscriptionsInRow,
		); err != nil {
			return nil, err
		}

		subscribers = append(subscribers, &subscriber)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscribers, nil
}

func (s *SubscriberStorage) Create(ctx context.Context, subscriber *model.Subscriber) error {
	const query = `
        INSERT INTO 
            subscribers
            (
             user_id,
             status,
             number_subscriptions, 
             subscription_time, 
             subscriptions_in_row
             )  
        VALUES ($1, $2, $3, $4, $5)  
        RETURNING 
        id
    `

	var id int
	err := s.db.QueryRowContext(
		ctx,
		query,
		subscriber.UserID,
		subscriber.StatusSubscription,
		subscriber.NumberSubscriptions,
		subscriber.SubscriptionTime,
		subscriber.SubscriptionsInRow,
	).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriberStorage) Update(ctx context.Context, subscriber *model.Subscriber, id int) error {
	const query = `
		UPDATE 
		    subscribers
		SET 
		    status_subscription = $1,
		    number_subscriptions = $2,
		    subscription_time = $3,
		    subscriptions_in_row = $4
		WHERE 
		    id = $5
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		subscriber.UserID,
		subscriber.StatusSubscription,
		subscriber.NumberSubscriptions,
		subscriber.SubscriptionTime,
		subscriber.SubscriptionsInRow,
		id,
	)

	return err
}

func (s *SubscriberStorage) Delete(ctx context.Context, id int) error {
	const query = `
		DELETE FROM subscribers WHERE id = $1
	`

	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
