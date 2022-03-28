package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type MessagesRepo struct {
	db *sql.DB
}

func NewMessagesRepo(db *sql.DB) *MessagesRepo {
	return &MessagesRepo{
		db: db,
	}
}

func (r *MessagesRepo) CreateMessage(inp *models.CreateMessage) (*model.NewMessage, error) {
	message := new(model.NewMessage)
	err := r.db.QueryRow(`
		INSERT INTO messages (room_id, msg_id, emp_id,target_id, body)
		SELECT $1, msg_count+1, $2, $3, $4
		FROM msg_state WHERE room_id = $1
		RETURNING msg_id, room_id, target_id, emp_id, body, created_at
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
	)

	return message, err
}
