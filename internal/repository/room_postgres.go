package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"
)

type RoomsRepo struct {
	db *sql.DB
}

func NewRoomsRepo(db *sql.DB) *RoomsRepo {
	return &RoomsRepo{
		db: db,
	}
}

// without filter
func (r *RoomsRepo) Rooms(employeeID int, params *model.Params) (*model.Rooms, error) {
	rooms := &model.Rooms{
		Rooms: []*model.Room{},
	}

	rows, err := r.db.Query(`
		SELECT r.room_id, r.name, r.view, m.emp_id, m.last_msg_read, c.last_msg_id
		FROM rooms r
		JOIN members m 
		    ON m.room_id = r.room_id AND m.emp_id = $1 
		JOIN msg_state c 
		    ON r.room_id = c.room_id
		LIMIT $6
		OFFSET $7
	`,
		employeeID,
		params.Limit,
		params.Offset,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := new(model.Room)
		if err = rows.Scan(&m.RoomID, &m.Name, &m.View, &m.LastMessageRead, &m.LastMessageID); err != nil {
			return nil, err
		}

		rooms.Rooms = append(rooms.Rooms, m)
	}

	return rooms, nil
}

func (r *RoomsRepo) FindRooms(employeeID int, inp *model.FindRooms, params *model.Params) (*model.Rooms, error) {
	rooms := &model.Rooms{
		Rooms: []*model.Room{},
	}
	if inp.Name != nil {
		*inp.Name = "%" + *inp.Name + "%"
	}

	rows, err := r.db.Query(`
		SELECT r.room_id, r.name, r.view, m.emp_id, m.last_msg_read, c.last_msg_id
		FROM rooms r
		JOIN members m 
		    ON m.room_id = r.room_id AND m.emp_id = $1 
		JOIN msg_state c 
		    ON r.room_id = c.room_id
		WHERE (
			    $2::BIGINT IS NULL 
			    OR r.room_id = $2 
			)
			AND (
			    $3::VARCHAR IS NULL 
			    OR r.name ILIKE $3
			)
		LIMIT $6
		OFFSET $7
	`,
		employeeID,
		inp.RoomID,
		inp.Name,
		params.Limit,
		params.Offset,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := new(model.Room)
		if err = rows.Scan(&m.RoomID, &m.Name, &m.View, &m.LastMessageRead, &m.LastMessageID); err != nil {
			return nil, err
		}

		rooms.Rooms = append(rooms.Rooms, m)
	}

	return rooms, nil
}
