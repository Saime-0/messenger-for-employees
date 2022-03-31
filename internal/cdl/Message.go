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
		MessageID int
	}
	messageResult struct {
		Message *model.Message
	}
)

func (d *Dataloader) Message(messageID int) (*model.Message, error) {
	res := <-d.categories.Message.addBaseRequest(
		&messageInp{
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
		messageIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		messageIDs = append(messageIDs, query.Inp.(*messageInp).MessageID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, 
		       coalesce(m.room_id, 0), 
		       coalesce(m.msg_id, 0), 
		       coalesce(m.emp_id, 0), 
		       coalesce(m.target_id, NULL), 
		       coalesce(m.body, ''), 
		       coalesce(m.created_at, 0)
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, messageid)
		LEFT JOIN messages m ON m.msg_id = inp.messageid
		`,
		pq.Array(ptrs),
		pq.Array(messageIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("message:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr      chanPtr
		targetID *int
	)
	for rows.Next() {
		m := &model.Message{Room: new(model.Room)}
		if err = rows.Scan(&ptr, &m.Room.RoomID, &m.MsgID, &m.Employee.EmpID, &targetID, &m.Body, &m.CreatedAt); err != nil {
			//c.Dataloader.healer.Alert("message (scan rows):" + err.Error())
			c.Error = err
			return
		}
		if m.MsgID == 0 {
			m = nil
		}
		if targetID != nil {
			m.TargetMsgID = &model.Message{MsgID: *targetID}
		}

		request := c.getRequest(ptr)
		request.Result.(*messageResult).Message = m
	}

	c.Error = nil
}
