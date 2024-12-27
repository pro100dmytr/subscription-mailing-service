package message

import (
	"context"
	"database/sql"
	"errors"
	"subscription-mailing-service/internal/config"
	"subscription-mailing-service/internal/model"
	"subscription-mailing-service/storage/postgres"
)

type MessageStorage struct {
	db *sql.DB
}

func (s *MessageStorage) Close() error {
	return postgres.CloseConnection(s.db)
}

func NewMessageStorage(cfg *config.Config) (*MessageStorage, error) {
	db, err := postgres.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &MessageStorage{db: db}, nil
}

func (s *MessageStorage) Get(ctx context.Context, id int) (*model.Message, error) {
	const query = `SELECT message FROM messages WHERE id = $1`
	message := &model.Message{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&message.Message,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	message.ID = id
	return message, err
}

func (s *MessageStorage) GetAll(ctx context.Context) ([]*model.Message, error) {
	const query = `SELECT id, message FROM messages`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		message := &model.Message{}
		if err := rows.Scan(
			&message.ID,
			&message.Message,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (s *MessageStorage) Create(ctx context.Context, message *model.Message) (*model.Message, error) {
	const query = `
       INSERT INTO messages (message)
       VALUES ($1)
       RETURNING id
   `

	var id int
	err := s.db.QueryRowContext(
		ctx,
		query,
		message.Message,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	message.ID = id
	return message, nil
}

func (s *MessageStorage) Update(ctx context.Context, message *model.Message, id int) error {
	const query = `UPDATE messages SET message = $1 WHERE id = $2`
	result, err := s.db.ExecContext(ctx, query, message.Message, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	message.ID = id

	return nil
}

func (s *MessageStorage) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM messages WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
