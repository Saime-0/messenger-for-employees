package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/models"
)

type PreparesRepo struct {
	db *sql.DB
}

func NewPreparesRepo(db *sql.DB) *PreparesRepo {
	return &PreparesRepo{
		db: db,
	}
}

func (r *PreparesRepo) ScheduleRefreshSessions(before int64) ([]*models.ScheduleRefreshSession, error) {
	var sessions []*models.ScheduleRefreshSession

	rows, err := r.db.Query(`
		SELECT id, expires_at
		FROM refresh_sessions
		WHERE $1 = 0 OR expires_at <= $1
		`,
		before,
	)
	if err != nil {
		return sessions, err
	}
	defer rows.Close()
	for rows.Next() {
		rs := &models.ScheduleRefreshSession{}
		if err = rows.Scan(&rs.ID, &rs.Exp); err != nil {
			return sessions, err
		}
		sessions = append(sessions, rs)
	}

	return sessions, nil
}
