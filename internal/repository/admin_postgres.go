package repository

import (
	"database/sql"
	"github.com/saime-0/messenger-for-employee/internal/models"
)

type AdminsRepo struct {
	db *sql.DB
}

func NewAdminsRepo(db *sql.DB) *AdminsRepo {
	return &AdminsRepo{
		db: db,
	}
}

func (r *AdminsRepo) AdminByToken(token string) (admin *models.Admin, err error) {
	err = r.db.QueryRow(`
		SELECT admin_id, email, token
		FROM admins
		WHERE token = $1`,
		token,
	).Scan(
		&admin.ID,
		&admin.Email,
		&admin.Token,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

func (r *AdminsRepo) AdminByID(id int) (admin *models.Admin, err error) {
	err = r.db.QueryRow(`
		SELECT admin_id, email, token
		FROM admins
		WHERE admin_id = $1`,
		id,
	).Scan(
		&admin.ID,
		&admin.Email,
		&admin.Token,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}
