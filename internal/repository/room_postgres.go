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

func (r *RoomsRepo) FindRooms(employeeID int, inp *model.FindRooms, params *model.Params) (*model.Rooms, error) {
	rooms := &model.Rooms{
		Rooms: []*model.Room{},
	}
	if inp.Name != nil {
		*inp.Name = "%" + *inp.Name + "%"
	}

	rows, err := r.db.Query(`
		SELECT array_position(e.room_seq, r.id), r.id, r.name, r.photo_url, r.view, m.emp_id, m.last_msg_read, c.last_msg_id, m.notify
		FROM rooms r
	    JOIN employees e
			ON e.id = $1
		JOIN members m
		    ON m.room_id = r.id AND m.emp_id = e.id
		JOIN msg_state c
		    ON r.id = c.room_id
		WHERE (
			    $2::BIGINT IS NULL
			    OR r.id = $2
			)
			AND (
			    $3::VARCHAR IS NULL
			    OR lower(r.name) ILIKE lower($3)
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
		if err = rows.Scan(&m.Pos, &m.RoomID, &m.Name, &m.PhotoURL, &m.View, &m.LastMessageRead, &m.LastMessageID, &m.Notify); err != nil {
			return nil, err
		}

		rooms.Rooms = append(rooms.Rooms, m)
	}

	return rooms, nil
}

func (r *RoomsRepo) FindMessages(inp *model.FindMessages, params *model.Params) (*model.Messages, error) {

	if inp.TextFragment != nil {
		*inp.TextFragment = "%" + *inp.TextFragment + "%"
	}

	var rows, err = r.db.Query(`
		SELECT m.id, room_id, emp_id, reply_id, body, created_at, prev, next
		FROM messages m
		WHERE (
		    $1::BIGINT IS NULL 
		    OR m.id = $1
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
		    OR m.reply_id = $4
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

	return parseMessagesRows(rows)
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
		INSERT INTO rooms (name, photo_url, view) VALUES ($1, $2)
		RETURNING id
	`,
		room.Name,
		room.PhotoUrl,
		room.View,
	).Scan(&roomID)
	return
}

func (r RoomsRepo) DropRoom(room *request_models.DropRoom) (err error) {
	_, err = r.db.Exec(`
		DELETE FROM rooms
		WHERE id = $1
	`,
		room.RoomID,
	)
	return
}

func (r RoomsRepo) RoomExists(roomID int) (exists bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1
		    FROm rooms
		    WHERE id = $1
		)
	`,
		roomID,
	).Scan(&exists)
	return
}

func (r RoomsRepo) EmployeeIsReceivesNotify(empID, roomID int) (enabled bool, err error) {
	err = r.db.QueryRow(`
	    SELECT notify
	    FROm members
	    WHERE emp_id = $1 AND room_id = $2
	`,
		empID,
		roomID,
	).Scan(&enabled)
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
	_, err = r.db.Exec(`
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
	)
	return
}

func (r RoomsRepo) KickEmployeesFromRoom(inp *request_models.KickEmployeesFromRooms) (err error) {
	_, err = r.db.Exec(`
		DELETE FROM members
		WHERE room_id = ANY($1) AND emp_id = ANY($2)
	`,
		pq.Array(inp.Rooms),
		pq.Array(inp.Employees),
	)
	return
}

func (r RoomsRepo) MoveRoom(empID, roomID int, prevRoomID *int) (err error) {
	_, err = r.db.Exec(`
		SELECT move_room_in_the_sequence($1, $2, $3)
	`,
		empID,
		roomID,
		prevRoomID,
	)
	return
}

func (r RoomsRepo) ReadMessage(empID, roomID, msgID int) (err error) {
	_, err = r.db.Exec(`
		UPDATE members
		SET last_msg_read = $3
		WHERE emp_id = $1 AND room_id = $2 
		  AND (
		      members.last_msg_read IS NULL OR 
		      last_msg_read < $3
	      )
	`,
		empID,
		roomID,
		msgID,
	)
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

func (r *RoomsRepo) SetNotify(empID int, roomID int, val bool) (err error) {
	_, err = r.db.Exec(`
		UPDATE members
		SET notify = $3
		WHERE emp_id = $1 AND room_id = $2
	`,
		empID,
		roomID,
		val,
	)
	return
}

func (r *RoomsRepo) CreateMessage(inp *models.CreateMessage) (*model.NewMessage, error) {
	message := new(model.NewMessage)
	err := r.db.QueryRow(`
		WITH send AS (
			INSERT INTO messages (room_id, emp_id, reply_id, body)
			VALUES ($1, $2, $3, $4)
			RETURNING id, room_id, reply_id, emp_id, body, created_at, prev
		)
		UPDATE members m
		SET last_msg_read = s.id
		FROM send s
		WHERE m.emp_id = $2 AND m.room_id = $1
		RETURNING s.id, s.room_id, s.reply_id, s.emp_id, s.body, s.created_at, s.prev
		`,
		inp.RoomID,
		inp.EmployeeID,
		inp.TargetMsgID,
		inp.Body,
	).Scan(
		&message.MsgID,
		&message.RoomID,
		&message.TargetMsgID,
		&message.EmployeeID,
		&message.Body,
		&message.CreatedAt,
		&message.Prev,
	)

	return message, err
}

func (r *RoomsRepo) RoomMessagesByRange(byRange *model.ByRange, limit int) (*model.Messages, error) {

	//// start must less than end // Nope
	//if byRange.Start > byRange.InDirection {
	//	byRange.Start, byRange.InDirection = byRange.InDirection, byRange.Start
	//}
	var rows, err = r.db.Query(`
		with lessToGreat(val) as (
		    VALUES ($2::BIGINT <= $3::BIGINT)
		)
		SELECT id, room_id, emp_id, reply_id, body, created_at, prev, next
		FROM lessToGreat, messages
		WHERE room_id = $1 
		  AND(
		      lessToGreat.val = TRUE AND id >= $2 AND id <= $3
		      OR id <= $2 AND id >= $3
	      )
		order by 
			case when not lessToGreat.val then id end desc,
			case when lessToGreat.val then id end asc
		limit $4
	`,
		byRange.RoomID,
		byRange.Start,
		byRange.InDirection,
		limit,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	return parseMessagesRows(rows)
}

func (r *RoomsRepo) RoomMessagesByCreated(byCreated *model.ByCreated) (*model.Messages, error) {
	// lessToGreat = true if created == after
	// lessToGreat = false if created == before
	lessToGreat := true
	if byCreated.Created == model.MsgCreatedBefore {
		lessToGreat = false
	}
	var rows, err = r.db.Query(`
		SELECT id, room_id, emp_id, reply_id, body, created_at, prev, next
		FROM messages
		WHERE room_id = $1 
			AND (
				$3 = false AND id <= $2
				OR
				$3 = true AND id >= $2
	    	)
		order by $3, 
                 case when $3 then id end asc,
                 case when not $3 then id end desc
		limit $4
	`,
		byCreated.RoomID,
		byCreated.StartMsg,
		lessToGreat,
		byCreated.Count,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	return parseMessagesRows(rows)
}

func parseMessagesRows(rows *sql.Rows) (*model.Messages, error) {
	messages := &model.Messages{
		Messages: []*model.Message{},
	}
	for rows.Next() {
		m := &model.Message{Room: new(model.Room)}
		var (
			targetID   *int
			employeeID *int
		)

		if err := rows.Scan(&m.MsgID, &m.Room.RoomID, &employeeID, &targetID, &m.Body, &m.CreatedAt, &m.Prev, &m.Next); err != nil {
			return nil, err
		}
		if targetID != nil {
			m.TargetMsg = &model.Message{MsgID: *targetID, Room: &model.Room{RoomID: m.Room.RoomID}}
		}
		if employeeID != nil {
			m.Employee = &model.Employee{EmpID: *employeeID}
		}

		messages.Messages = append(messages.Messages, m)
	}
	return messages, nil
}
