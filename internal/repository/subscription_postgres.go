package repository

import (
	"database/sql"
	"github.com/lib/pq"
)

type QueryUserGroup func(objectIDs ...int) (users []int, err error)

type SubscribersRepo struct {
	db          *sql.DB
	Members     QueryUserGroup
	RoomReaders QueryUserGroup
}

func NewSubscribersRepo(db *sql.DB) *SubscribersRepo {
	sub := &SubscribersRepo{
		db: db,
	}
	sub.initFuncs()
	return sub
}

func completeIntArray(rows *sql.Rows) (arr []int, err error) {
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		arr = append(arr, id)
	}
	return arr, nil
}

func (r *SubscribersRepo) initFuncs() {

	r.Members = func(roomIDs ...int) (employees []int, err error) {
		rows, err := r.db.Query(`
			SELECT emp_id 
			FROM members
			WHERE room_id = ANY ($1)`,
			pq.Array(roomIDs),
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		employees, err = completeIntArray(rows)
		if err != nil {
			return nil, err
		}

		return employees, nil
	}

}
