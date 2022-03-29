package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
)

type ChatsRepo struct {
	db *sql.DB
}

func NewChatsRepo(db *sql.DB) *ChatsRepo {
	return &ChatsRepo{
		db: db,
	}
}

func (r *ChatsRepo) Tags(params *model.Params) (*model.Tags, error) {

	tags := &model.Tags{
		Tags: []*model.Tag{},
	}

	var rows, err = r.db.Query(`
		SELECT t.tag_id, t.name
		FROM tags t
		LIMIT $1
		OFFSET $2
		`,
		params.Limit,
		params.Offset,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := new(model.Tag)
		if err = rows.Scan(&m.TagID, &m.Name); err != nil {
			return nil, err
		}
		tags.Tags = append(tags.Tags, m)
	}

	return tags, nil
}

func (r *ChatsRepo) FindMessages(empID int, inp *model.FindMessages, params *model.Params) (*model.Messages, error) {

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
func (r *ChatsRepo) UserHasAccessToRooms(employeeID int, chats []int) (noAccessTo int, err error) {
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
