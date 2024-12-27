package mail

import (
	"context"
	"database/sql"
	"errors"
	"subscription-mailing-service/internal/config"
	"subscription-mailing-service/internal/model"
	"subscription-mailing-service/storage/postgres"
	"time"
)

type MailStorage struct {
	db *sql.DB
}

func (s *MailStorage) Close() error {
	return postgres.CloseConnection(s.db)
}

func NewMailStorage(cfg *config.Config) (*MailStorage, error) {
	db, err := postgres.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &MailStorage{db: db}, nil
}

func (s *MailStorage) Get(ctx context.Context, id int) (*model.Mail, error) {
	const query = `SELECT to_list, subject, body, content_type, sent_at FROM mails WHERE id = $1`
	mail := &model.Mail{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&mail.To,
		&mail.Subject,
		&mail.Body,
		&mail.ContentType,
		&mail.SentAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	mail.ID = id
	return mail, nil
}

func (s *MailStorage) GetAll(ctx context.Context) ([]*model.Mail, error) {
	const query = `SELECT to_list, subject, body, content_type, sent_at FROM mails`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mails := []*model.Mail{}
	for rows.Next() {
		mail := &model.Mail{}
		if err := rows.Scan(
			&mail.To,
			&mail.Subject,
			&mail.Body,
			&mail.ContentType,
			&mail.SentAt,
		); err != nil {
			return nil, err
		}
		mails = append(mails, mail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mails, rows.Err()
}

func (s *MailStorage) Create(ctx context.Context, mail *model.Mail) (*model.Mail, error) {
	const query = `
		INSERT INTO mails(to_list, subject, body, content_type, sent_at)
		VALUES ($1, $2, $3, $4, $5) 
		        RETURNING id
		`

	var id int
	err := s.db.QueryRowContext(ctx, query, mail.To, mail.Subject, mail.Body, mail.ContentType, time.Now()).Scan(&id)
	if err != nil {
		return nil, err
	}

	mail.ID = id

	return mail, nil
}

// TODO сделать отправку
//func (s *MailStorage) Send() {
//
//}

func (s *MailStorage) Update(ctx context.Context, mail *model.Mail, id int) error {
	const query = `
		UPDATE
		    mails
		SET 
		    to_list = $1,
		    subject = $2,
		    body = $3,
		    content_type = $4, 
		    sent_at = $5
		WHERE 
		    id =$6`

	result, err := s.db.ExecContext(ctx, query, mail.To, mail.Subject, mail.Body, mail.ContentType, mail.SentAt, id)
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

func (s *MailStorage) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM mails WHERE id = $1`
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

//func (s *MailStorage) SearchMails(ctx context.Context, query string) ([]*model.Mail, error) {
//	const sql = `SELECT * FROM mails WHERE subject ILIKE $1 OR body ILIKE $1`
//	rows, err := s.db.QueryContext(ctx, sql, "%"+query+"%")
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var mails []*model.Mail
//	for rows.Next() {
//		mail := &model.Mail{}
//		if err := rows.Scan(&mail.To, &mail.Subject, &mail.Body, &mail.ContentType, &mail.SentAt); err != nil {
//			return nil, err
//		}
//		mails = append(mails, mail)
//	}
//	return mails, rows.Err()
//}
