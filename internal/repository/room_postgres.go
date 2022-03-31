package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/admin/request_models"
	"github.com/saime-0/messenger-for-employee/internal/models"
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
func (r *RoomsRepo) EmployeeRooms(employeeID int, params *model.Params) (*model.Rooms, error) {
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

func (r *RoomsRepo) FindMessages(empID int, inp *model.FindMessages, params *model.Params) (*model.Messages, error) {

	messages := &model.Messages{
		Messages: []*model.Message{},
	}
	if inp.TextFragment != nil {
		*inp.TextFragment = "%" + *inp.TextFragment + "%"
	}

	var rows, err = r.db.Query(`
		SELECT m.room_id, m.msg_id, m.emp_id, m.target_id, m.body, m.created_at
		FROM messages m
		WHERE (
		    $1::BIGINT IS NULL 
		    OR m.msg_id = $1
		)
		AND (
		    $2::BIGINT IS NULL 
		    OR m.emp_id = $2 
		)
		AND (
		    $3::BIGINT IS NULL 
		    OR m.room_id = $3
		)
		AND (
		    $4::BIGINT IS NULL 
		    OR m.target_id = $4
		)
		AND (
		    $5::VARCHAR IS NULL 
		    OR body ILIKE $5
		)
		LIMIT $6
		OFFSET $7
		`,
		inp.MsgID,
		inp.EmpID,
		inp.RoomID,
		inp.TargetID,
		inp.TextFragment,
		params.Limit,
		params.Offset,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := &model.Message{
			Room:     new(model.Room),
			Employee: new(model.Employee),
		}
		var (
			targetID *int
		)
		if err = rows.Scan(&m.Room.RoomID, &m.MsgID, &m.Employee.EmpID, &targetID, &m.Body, &m.CreatedAt); err != nil {
			return nil, err
		}
		if targetID != nil {
			m.TargetMsgID = &model.Message{MsgID: *targetID}
		}
		messages.Messages = append(messages.Messages, m)
	}

	return messages, nil
}

// if noAccessTo = 0 then acces allow to all chats
func (r *RoomsRepo) EmployeeHasAccessToRooms(employeeID int, chats []int) (noAccessTo int, err error) {
	rows, err := r.db.Query(`
	    SELECT roomid, m.room_id is not null as is_member
	    FROM unnest($2::BIGINT[]) roomid
	    LEFT JOIN members m ON m.emp_id = $1 AND roomid = m.room_id`,
		employeeID,
		pq.Array(chats),
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		sr := new(models.SubRoom)
		if err = rows.Scan(&sr.RoomID, &sr.IsMember); err != nil {
			return
		}
		if !sr.IsMember {
			return *sr.RoomID, nil
		}
	}
	return 0, nil
}

func (r RoomsRepo) CreateRoom(room *request_models.CreateRoom) (roomID int, err error) {
	err = r.db.QueryRow(`
		INSERT INTO rooms (name, view) VALUES ($1, $2)
		RETURNING room_id
	`,
		room.Name,
		room.View,
	).Scan(&roomID)
	return
}

func (r RoomsRepo) DropRoom(room *request_models.DropRoom) (err error) {
	err = r.db.QueryRow(`
		DELETE FROM rooms
		WHERE room_id = $1
	`,
		room.RoomID,
	).Err()
	return
}

func (r RoomsRepo) EmployeesIsNotMember(roomID int, empIDs ...int) (empIsMember int, err error) {
	err = r.db.QueryRow(`
		SELECT coalesce((
	        SELECT emp_id
			FROM members
		    WHERE room_id = $1 AND emp_id = ANY($2)
	        LIMIT 1
		), 0)
	`,
		roomID,
		pq.Array(empIDs),
	).Scan(&empIsMember)
	return
}

func (r RoomsRepo) AddEmployeeToRoom(inp *request_models.AddEmployeeToRooms) (err error) {
	err = r.db.QueryRow(`
			WITH "except"(room_id) AS (
			    SELECT room_id
				FROM members
			    WHERE emp_id = $1 AND room_id = ANY($2)
			)
			INSERT INTO members (emp_id, room_id) 
			SELECT $1, roomid
			FROM unnest($2::bigint[]) inp(roomid)
			WHERE roomid != ALL(select room_id from "except")
	`,
		inp.Employee,
		pq.Array(inp.Rooms),
	).Err()
	return
}

func (r RoomsRepo) KickEmployeesFromRoom(inp *request_models.AddOrDeleteEmployeesInRoom) (err error) {
	err = r.db.QueryRow(`
		DELETE FROM members
		WHERE room_id = ANY($1) AND emp_id = ANY($2)
	`,
		pq.Array(inp.Rooms),
		pq.Array(inp.Employees),
	).Err()
	return
}

func (r RoomsRepo) RoomMembersID(roomID int) (employeeIDs []int, err error) {
	rows, err := r.db.Query(`
		SELECT emp_id
		FROM members
		WHERE room_id = $1
	`,
		roomID,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var empID int
		if err = rows.Scan(&empID); err != nil {
			return nil, err
		}

		employeeIDs = append(employeeIDs, empID)
	}

	return
}
