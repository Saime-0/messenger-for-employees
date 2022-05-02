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
// fix: перевернут порядок комнат
// unused
func (r *RoomsRepo) EmployeeRooms(employeeID int, params *model.Params) (*model.Rooms, error) {
	rooms := &model.Rooms{
		Rooms: []*model.Room{},
	}

	rows, err := r.db.Query(`
		WITH seq(seq) AS (
		    SELECT room_seq[
		        (select 1+coalesce($2,0)):
                (select 1+coalesce($2,0)) + (select coalesce($3, array_length(room_seq, 1)))] 
		    FROM employees WHERE emp_id = $1
		)
		SELECT r.room_id, r.name, r.view, m.emp_id, m.last_msg_read, c.last_msg_id, m.prev_id
		FROM rooms r
		JOIN seq 
		    ON r.room_id = ANY(seq.seq)
		JOIN members m 
		    ON m.room_id = r.room_id AND m.emp_id = $1 
		JOIN msg_state c 
		    ON r.room_id = c.room_id
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
		if err = rows.Scan(&m.RoomID, &m.Name, &m.View, &m.LastMessageRead, &m.LastMessageID, &m.PrevRoomID); err != nil {
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
		SELECT r.room_id, r.name, r.view, m.emp_id, m.last_msg_read, c.last_msg_id, m.prev_id
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
		if err = rows.Scan(&m.RoomID, &m.Name, &m.View, &m.LastMessageRead, &m.LastMessageID, &m.PrevRoomID); err != nil {
			return nil, err
		}

		rooms.Rooms = append(rooms.Rooms, m)
	}

	return rooms, nil
}

func (r *RoomsRepo) FindMessages(inp *model.FindMessages, params *model.Params) (*model.Messages, error) {

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
			m.TargetMsg = &model.Message{MsgID: *targetID, Room: &model.Room{RoomID: m.Room.RoomID}}
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

func (r RoomsRepo) RoomExists(roomID int) (exists bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1
		    FROm rooms
		    WHERE room_id = $1
		)
	`,
		roomID,
	).Scan(&exists)
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
		inp.EmpID,
		pq.Array(inp.Rooms),
	).Err()
	return
}

func (r RoomsRepo) KickEmployeesFromRoom(inp *request_models.KickEmployeesFromRooms) (err error) {
	err = r.db.QueryRow(`
		DELETE FROM members
		WHERE room_id = ANY($1) AND emp_id = ANY($2)
	`,
		pq.Array(inp.Rooms),
		pq.Array(inp.Employees),
	).Err()
	return
}

func (r RoomsRepo) MoveRoom(empID, roomID int, prevRoomID *int) (err error) {
	err = r.db.QueryRow(`
		SELECT move_room_in_the_sequence($1, $2, $3)
	`,
		empID,
		roomID,
		prevRoomID,
	).Err()
	return
}

func (r RoomsRepo) ReadMessage(empID, roomID, msgID int) (err error) {
	err = r.db.QueryRow(`
		UPDATE members
		SET last_msg_read = $3
		WHERE emp_id = $1 AND room_id = $2 AND last_msg_read < $3
	`,
		empID,
		roomID,
		msgID,
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

func (r *RoomsRepo) CreateMessage(inp *models.CreateMessage) (*model.NewMessage, error) {
	message := new(model.NewMessage)
	err := r.db.QueryRow(`
		INSERT INTO messages (room_id, msg_id, emp_id,target_id, body)
		SELECT $1, msg_count+1, $2, $3, $4
		FROM msg_state WHERE room_id = $1
		RETURNING msg_id, room_id, target_id, emp_id, body, created_at, prev
		`,
		inp.RoomID,
		inp.EmployeeID,
		inp.TargetMsgID,
		inp.Body,
	).Scan(
		&message.MsgID,
		&message.RoomID,
		&message.TargetMsgID,
		&message.EmpID,
		&message.Body,
		&message.CreatedAt,
		&message.Prev,
	)

	return message, err
}

func (r *RoomsRepo) RoomMessagesByRange(byRange *model.ByRange, limit int) (*model.Messages, error) {
	messages := &model.Messages{
		Messages: []*model.Message{},
	}
	// start must less than end
	if byRange.Start > byRange.End {
		byRange.Start, byRange.End = byRange.End, byRange.Start
	}
	var rows, err = r.db.Query(`
		SELECT room_id, msg_id, emp_id, target_id, body, created_at, prev, next
		FROM messages
		WHERE room_id = $1 AND msg_id >= $2 AND msg_id <= $3
		limit $4
	`,
		byRange.RoomID,
		byRange.Start,
		byRange.End,
		limit,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := &model.Message{Room: new(model.Room), Employee: new(model.Employee)}
		var targetID *int

		if err = rows.Scan(&m.Room.RoomID, &m.MsgID, &m.Employee.EmpID, &targetID, &m.Body, &m.CreatedAt, &m.Prev, &m.Next); err != nil {
			return nil, err
		}
		if targetID != nil {
			m.TargetMsg = &model.Message{MsgID: *targetID, Room: &model.Room{RoomID: m.Room.RoomID}}
		}

		messages.Messages = append(messages.Messages, m)
	}

	return messages, nil
}

func (r *RoomsRepo) RoomMessagesByCreated(byCreated *model.ByCreated) (*model.Messages, error) {
	messages := &model.Messages{
		Messages: []*model.Message{},
	}
	// lessToGreat = true if created == after
	// lessToGreat = false if created == before
	lessToGreat := true
	if byCreated.Created == model.MsgCreatedBefore {
		lessToGreat = false
	}
	var rows, err = r.db.Query(`
		SELECT room_id, msg_id, emp_id, target_id, body, created_at, prev, next
		FROM messages
		WHERE room_id = $1 
			AND (
				$3 = false AND msg_id <= $2
				OR
				$3 = true AND msg_id >= $2
	    	)
		order by $3, 
                 case when $3 then msg_id end asc,
                 case when not $3 then msg_id end desc
		limit $4
	`,
		byCreated.RoomID,
		byCreated.StartMsg,
		lessToGreat,
		byCreated.Count,
	)
	defer rows.Close()
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	for rows.Next() {
		m := &model.Message{Room: new(model.Room), Employee: new(model.Employee)}
		var targetID *int

		if err = rows.Scan(&m.Room.RoomID, &m.MsgID, &m.Employee.EmpID, &targetID, &m.Body, &m.CreatedAt, &m.Prev, &m.Next); err != nil {
			return nil, err
		}
		if targetID != nil {
			m.TargetMsg = &model.Message{MsgID: *targetID, Room: &model.Room{RoomID: m.Room.RoomID}}
		}

		messages.Messages = append(messages.Messages, m)
	}

	return messages, nil
}
