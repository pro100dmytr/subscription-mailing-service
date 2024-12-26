package mail

import (
	"database/sql"
	"subscription-mailing-service/internal/config"
	"subscription-mailing-service/storage/postgres"
)

type MailStorage struct {
	db *sql.DB
}

func (s *MailStorage) Close() error {
	return postgres.CloseConnection(s.db)
}

func NewStorage(cfg *config.Config) (*MailStorage, error) {
	db, err := postgres.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &MailStorage{db: db}, nil
}

func (s *MailStorage) GetMailInfo() {

}

func (s *MailStorage) GetAllMails() {

}

func (s *MailStorage) SendMail() {

}

func (s *MailStorage) UpdateMail() {

}

func (s *MailStorage) DeleteMail() {

}
