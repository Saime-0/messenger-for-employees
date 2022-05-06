package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/messenger-for-employee/graph/model"
)

func (r *roomResult) isRequestResult() {}
func (r *roomInp) isRequestInput()     {}

type (
	roomInp struct {
		EmpID  int
		RoomID int
	}
	roomResult struct {
		Room *model.Room
	}
)

func (d *Dataloader) Room(empID, roomID int) (*model.Room, error) {
	res := <-d.categories.Room.addBaseRequest(
		&roomInp{
			EmpID:  empID,
			RoomID: roomID,
		},
		new(roomResult),
	)
	if res == nil {
		return nil, d.categories.Room.Error
	}
	return res.(*roomResult).Room, nil
}

func (c *parentCategory) room() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		empIDs  []int
		roomIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		empIDs = append(empIDs, query.Inp.(*roomInp).EmpID)
		roomIDs = append(roomIDs, query.Inp.(*roomInp).RoomID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr,
		       array_position(e.room_seq, r.id),
				coalesce(r.id, 0),
				coalesce(r.name, ''),
				coalesce(r.view, 'TALK'),
				m.last_msg_read,
				c.last_msg_id
		FROM unnest($1::varchar[], $2::bigint[], $3::bigint[]) inp(ptr, empid, roomid)
		LEFT JOIN employees e
		    ON e.id = inp.empid
		LEFT JOIN members m 
		    ON m.emp_id = e.id AND m.room_id = inp.roomid
		LEFT JOIN rooms r ON r.id = m.room_id
		LEFT JOIN msg_state c ON c.room_id = m.room_id
		`,
		pq.Array(ptrs),
		pq.Array(empIDs),
		pq.Array(roomIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("room:" + err.Desk())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // Каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := new(model.Room)
		if err = rows.Scan(&ptr, &m.Pos, &m.RoomID, &m.Name, &m.View, &m.LastMessageRead, &m.LastMessageID); err != nil {
			//c.Dataloader.healer.Alert("room (scan rows):" + err.Desk())
			c.Error = err
			return
		}

		if m.RoomID == 0 {
			m = nil
		}
		request := c.getRequest(ptr)
		request.Result.(*roomResult).Room = m
	}

	c.Error = nil
}
