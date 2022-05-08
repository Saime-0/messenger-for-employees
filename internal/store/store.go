package store

import (
	"database/sql"
	"github.com/saime-0/messenger-for-employee/internal/config"
)

func InitDB(cfg *config.Config2) (*sql.DB, error) {
	// connection string

	// open database
	db, err := sql.Open("postgres", cfg.PostgresConnection)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// pq: sorry, too many clients already fix 0_0 maybe... and replace QueryRow().Err() -> Exec()
	db.SetMaxIdleConns(0)
	db.SetConnMaxIdleTime(0)
	return db, nil

}
