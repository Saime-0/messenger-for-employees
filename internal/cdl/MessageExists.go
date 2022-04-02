package cdl

import (
	"fmt"
	"github.com/lib/pq"
)

func (r *messageExistsResult) isRequestResult() {}
func (r *messageExistsInp) isRequestInput()     {}

type (
	messageExistsInp struct {
		RoomID int
		MsgID  int
	}
	messageExistsResult struct {
		Exists bool
	}
)

func (d *Dataloader) MessageExists(roomID, msgID int) (bool, error) {
	d.healer.Debug("Dataloader: новый запрос MessageExists")
	res := <-d.categories.MessageExists.addBaseRequest(
		&messageExistsInp{
			RoomID: roomID,
			MsgID:  msgID,
		},
		new(messageExistsResult),
	)
	if res == nil {
		return false, d.categories.MessageExists.Error
	}
	return res.(*messageExistsResult).Exists, nil
}

func (c *parentCategory) messageExists() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		roomIDs []int
		msgIDs  []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		roomIDs = append(roomIDs, query.Inp.(*messageExistsInp).RoomID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, msg_id is not null
		FROM unnest($1::varchar[], $2::bigint[], $3::bigint[]) inp(ptr, roomid, msgid)
		LEFT JOIN messages m ON m.room_id = inp.roomid AND m.msg_id = inp.msgid
		`,
		pq.Array(ptrs),
		pq.Array(roomIDs),
		pq.Array(msgIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("messageExists:" + err.Desk())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr    chanPtr
		exists bool
	)
	for rows.Next() {

		if err = rows.Scan(&ptr, &exists); err != nil {
			//c.Dataloader.healer.Alert("messageExists (scan rows):" + err.Desk())
			c.Error = err
			return
		}

		request := c.getRequest(ptr)
		request.Result.(*messageExistsResult).Exists = exists
	}

	c.Error = nil
}
