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

func NewSubscriberStorage(cfg *config.Config) (*SubscriberStorage, error) {
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
		    subscriptions_in_row,
		    subscriptions_level
		    
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
		&subscriber.SubscriptionLevel,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	subscriber.ID = id
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
		    subscriptions_in_row,
		    subscriptions_level
		FROM
		    subscribers
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscribers := []*model.Subscriber{}
	for rows.Next() {
		subscriber := &model.Subscriber{}
		if err := rows.Scan(
			&subscriber.ID,
			&subscriber.UserID,
			&subscriber.StatusSubscription,
			&subscriber.NumberSubscriptions,
			&subscriber.SubscriptionTime,
			&subscriber.SubscriptionsInRow,
			&subscriber.SubscriptionLevel,
		); err != nil {
			return nil, err
		}

		subscribers = append(subscribers, subscriber)
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
            status_subscription,
            number_subscriptions,
            subscription_time,
            subscriptions_in_row,
            subscriptions_level
            )
       VALUES ($1, $2, $3, $4, $5, $6)
       RETURNING id
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
		subscriber.SubscriptionLevel,
	).Scan(&id)

	if err != nil {
		return err
	}

	subscriber.ID = id

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
		subscriber.StatusSubscription,
		subscriber.NumberSubscriptions,
		subscriber.SubscriptionTime,
		subscriber.SubscriptionsInRow,
		id,
	)

	subscriber.ID = id

	return err
}

func (s *SubscriberStorage) Delete(ctx context.Context, id int) error {
	const query = `
		DELETE FROM subscribers WHERE id = $1
	`

	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *SubscriberStorage) LevelUp(ctx context.Context, subscriber *model.Subscriber, id int) error {
	const query = `UPDATE subscribers SET subscriptions_level = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, subscriber.SubscriptionLevel, id)

	return err
}

func (s *SubscriberStorage) GetByLevel(ctx context.Context, level string) ([]*model.Subscriber, error) {
	const query = `
		SELECT
		    id,
		    user_id,
		    status_subscription,
		    number_subscriptions,
		    subscription_time,
		    subscriptions_in_row,
		    subscriptions_level
       FROM
           subscribers
       WHERE
           subscriptions_level = $1
           `

	rows, err := s.db.QueryContext(ctx, query, level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribers []*model.Subscriber
	for rows.Next() {
		subscriber := &model.Subscriber{}
		if err := rows.Scan(
			&subscriber.ID,
			&subscriber.UserID,
			&subscriber.StatusSubscription,
			&subscriber.NumberSubscriptions,
			&subscriber.SubscriptionTime,
			&subscriber.SubscriptionsInRow,
			&subscriber.SubscriptionLevel,
		); err != nil {
			return nil, err
		}
		subscribers = append(subscribers, subscriber)
	}

	return subscribers, rows.Err()
}
