package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *RoomsResult) isRequestResult() {}
func (r *RoomsInp) isRequestInput()     {}

type (
	RoomsResult struct {
		Rooms *model.Rooms
	}
	RoomsInp struct {
		EmployeeID int
	}
)

func (d *Dataloader) Rooms(empID int) (*model.Rooms, error) {
	res := <-d.categories.Rooms.addBaseRequest(
		&RoomsInp{
			EmployeeID: empID,
		},
		&RoomsResult{
			Rooms: &model.Rooms{
				Rooms: []*model.Room{},
			},
		},
	)
	if res == nil {
		return nil, d.categories.Rooms.Error
	}
	return res.(*RoomsResult).Rooms, nil
}

func (c *parentCategory) rooms() {
	var (
		inp = c.Requests

		ptrs   []chanPtr
		empIDs []int
	)
	for _, query := range inp {
		empIDs = append(empIDs, query.Inp.(*RoomsInp).EmployeeID)
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr,
				coalesce(r.room_id, 0),
				coalesce(r.name, ''),
				coalesce(r.view, 'TALK'),
				coalesce(m.last_msg_read, 0),
				coalesce(c.last_msg_id, 0)
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, empid)
		LEFT JOIN members m ON m.emp_id = inp.empid
		LEFT JOIN rooms r on r.room_id = m.room_id
		LEFT JOIN msg_state c on r.room_id = c.room_id
		`,
		pq.Array(ptrs),
		pq.Array(empIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("rooms:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // Каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := new(model.Room)

		if err = rows.Scan(&ptr, &m.RoomID, &m.Name, &m.View, &m.LastMessageRead, &m.LastMessageID); err != nil {
			//c.Dataloader.healer.Alert("rooms (scan rows):" + err.Error())
			c.Error = err
			return
		}
		if m.RoomID == 0 {
			continue
		}
		request := c.getRequest(ptr)
		request.Result.(*RoomsResult).Rooms.Rooms = append(request.Result.(*RoomsResult).Rooms.Rooms, m)
	}

	c.Error = nil
}
