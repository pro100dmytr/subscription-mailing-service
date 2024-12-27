package db

import (
	"database/sql"
	"fmt"
)

const initTableUsersSQL = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    login VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255),
    password VARCHAR(255)
);`

const initTableSubscribersSQL = `
CREATE TABLE IF NOT EXISTS subscribers (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    status_subscription VARCHAR(255),
    number_subscriptions INT,
    subscription_time TIMESTAMP,
    subscriptions_in_row INT
);`

const initTableMessagesSQL = `
CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	message VARCHAR(255)
);`

const initTableMailsSQL = `
CREATE TABLE IF NOT EXISTS mails (
    id SERIAL PRIMARY KEY,         
    to_list TEXT[],             
    subject VARCHAR(255) NOT NULL,    
    body TEXT NOT NULL,              
    content_type VARCHAR(50),         
    sent_at TIMESTAMP 
);`

func InitDatabase(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("Failed to connect to database: %w", err)
	}

	_, err := db.Exec(initTableUsersSQL)
	if err != nil {
		return fmt.Errorf("Error creating user table: %w", err)
	}

	_, err = db.Exec(initTableSubscribersSQL)
	if err != nil {
		return fmt.Errorf("Error creating subscriber table: %w", err)
	}

	_, err = db.Exec(initTableMessagesSQL)
	if err != nil {
		return fmt.Errorf("Error creating message table: %w", err)
	}

	_, err = db.Exec(initTableMailsSQL)
	if err != nil {
		return fmt.Errorf("Error creating mail table: %w", err)
	}
	return nil
}
