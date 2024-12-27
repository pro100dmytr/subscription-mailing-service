package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	dbInit "subscription-mailing-service/db"
	"subscription-mailing-service/internal/config"
)

func OpenConnection(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.DBName,
		cfg.Database.Password,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = dbInit.InitDatabase(db)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CloseConnection(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}
