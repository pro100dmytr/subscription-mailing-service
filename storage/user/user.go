package user

import (
	"context"
	"database/sql"
	"errors"
	"subscription-mailing-service/internal/config"
	"subscription-mailing-service/internal/model"
	"subscription-mailing-service/storage/postgres"
)

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Close() error {
	return postgres.CloseConnection(s.db)
}

func NewStorage(cfg *config.Config) (*UserStorage, error) {
	db, err := postgres.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &UserStorage{db: db}, nil
}

func (s *UserStorage) Create(ctx context.Context, user *model.User) (*model.User, error) {
	const query = `
        INSERT INTO users (first_name, last_name, login, email, password)  
        VALUES ($1, $2, $3, $4, $5)  
        RETURNING id
    `

	var id int
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Login,
		user.Email,
		user.Password,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (s *UserStorage) Get(ctx context.Context, id int) (*model.User, error) {
	const query = `SELECT first_name, last_name, login, email, password FROM users WHERE id = $1`
	user := &model.User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Login,
		&user.Email,
		&user.Password,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

func (s *UserStorage) GetAll(ctx context.Context) ([]*model.User, error) {
	const query = `SELECT id, first_name, last_name, login, email, password FROM users`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Login,
			&user.Email,
			&user.Password,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (s *UserStorage) Update(ctx context.Context, user *model.User, id int) error {
	const query = `UPDATE users SET first_name = $1, last_name = $2, login = $3, email = $4, password = $5 WHERE id = $6`
	result, err := s.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Login, user.Email, user.Password, id)
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

func (s *UserStorage) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM users WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
