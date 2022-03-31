package repository

import (
	"database/sql"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/rules"

	"github.com/saime-0/messenger-for-employee/internal/models"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) CreateRefreshSession(employeeID int, sessionModel *models.RefreshSession, overflowDelete bool) (id int, err error) {
	err = r.db.QueryRow(`
		INSERT INTO refresh_sessions (emp_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id`,
		employeeID,
		sessionModel.RefreshToken,
		//sessionModel.UserAgent,
		sessionModel.ExpAt,
	).Scan(
		&id,
	)
	if err != nil {
		return
	}

	if overflowDelete {
		err = r.OverflowDelete(employeeID, rules.MaxRefreshSession)
		if err != nil {
			return
		}
	}

	return
}

func (r *AuthRepo) UpdateRefreshSession(sessionID int, sessionModel *models.RefreshSession) (err error) {
	err = r.db.QueryRow(`
		INSERT INTO refresh_sessions (emp_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id`,
		sessionID,
		sessionModel.RefreshToken,
		//sessionModel.UserAgent,
		sessionModel.ExpAt,
	).Err()

	return
}

func (r *AuthRepo) OverflowDelete(employeeID, limit int) (err error) {
	err = r.db.QueryRow(`
		DELETE FROM refresh_sessions 
		WHERE id IN(                 
		    WITH session_count AS (
		        SELECT count(1) AS val
				FROM refresh_sessions
				WHERE emp_id = $1
		        GROUP BY emp_id
		    )
		    SELECT id
		    FROM refresh_sessions 
		    WHERE coalesce((select val from session_count) > $2, false) = true AND emp_id = $1
		    ORDER BY id ASC 
		    LIMIT abs((select val from session_count) - $2)
		    
		    )`,
		employeeID,
		limit,
	).Err()
	if err != nil {
		return cerrors.Wrap(err, "не удалось удалить лишние сессии")
	}

	return
}

func (r *AuthRepo) FindSessionByComparedToken(token string) (sessionId int, employeeID int, err error) {
	err = r.db.QueryRow(`
		SELECT id, emp_id
		FROM refresh_sessions
		WHERE refresh_token = $1`,
		token,
	).Scan(
		&sessionId,
		&employeeID,
	)

	return
}
