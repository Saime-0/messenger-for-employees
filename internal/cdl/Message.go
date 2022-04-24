package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/messenger-for-employee/graph/model"
)

func (r *messageResult) isRequestResult() {}
func (r *messageInp) isRequestInput()     {}

type (
	messageInp struct {
		RoomID    int
		MessageID int
	}
	messageResult struct {
		Message *model.Message
	}
)

func (d *Dataloader) Message(roomID, messageID int) (*model.Message, error) {
	res := <-d.categories.Message.addBaseRequest(
		&messageInp{
			RoomID:    roomID,
			MessageID: messageID,
		},
		new(messageResult),
	)
	if res == nil {
		return nil, d.categories.Message.Error
	}
	return res.(*messageResult).Message, nil
}

func (c *parentCategory) message() {
	var (
		inp = c.Requests

		ptrs       []chanPtr
		roomIDs    []int
		messageIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		roomIDs = append(roomIDs, query.Inp.(*messageInp).RoomID)
		messageIDs = append(messageIDs, query.Inp.(*messageInp).MessageID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, 
		       coalesce(m.room_id, 0), 
		       coalesce(m.msg_id, 0), 
		       coalesce(m.emp_id, 0), 
		       m.target_id, 
		       coalesce(m.body, ''), 
		       coalesce(m.created_at, 0),
			   m.prev,
			   m.next
		FROM unnest($1::varchar[], $2::bigint[], $3::bigint[]) inp(ptr, roomid, messageid)
		LEFT JOIN messages m ON m.msg_id = inp.messageid AND m.room_id = inp.roomid
		`,
		pq.Array(ptrs),
		pq.Array(roomIDs),
		pq.Array(messageIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("message:" + err.Desk())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // Каждую итерацию будем менять значения
		ptr      chanPtr
		targetID *int
	)
	for rows.Next() {
		m := &model.Message{Room: new(model.Room), Employee: new(model.Employee)}
		if err = rows.Scan(&ptr, &m.Room.RoomID, &m.MsgID, &m.Employee.EmpID, &targetID, &m.Body, &m.CreatedAt, &m.Prev, &m.Next); err != nil {
			//c.Dataloader.healer.Alert("message (scan rows):" + err.Desk())
			c.Error = err
			return
		}
		if m.MsgID == 0 {
			m = nil
		}
		if targetID != nil {
			m.TargetMsg = &model.Message{MsgID: *targetID, Room: &model.Room{RoomID: m.Room.RoomID}}
		}

		request := c.getRequest(ptr)
		request.Result.(*messageResult).Message = m
	}

	c.Error = nil
}
