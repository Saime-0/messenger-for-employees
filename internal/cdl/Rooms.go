package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/messenger-for-employee/graph/model"
)

func (r *RoomsResult) isRequestResult() {}
func (r *RoomsInp) isRequestInput()     {}

type (
	RoomsResult struct {
		Rooms *model.Rooms
	}
	RoomsInp struct {
		EmployeeID int
		Params     *model.Params
	}
)

func (d *Dataloader) Rooms(empID int, params *model.Params) (*model.Rooms, error) {
	res := <-d.categories.Rooms.addBaseRequest(
		&RoomsInp{
			EmployeeID: empID,
			Params:     params,
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

		ptrs    []chanPtr
		empIDs  []int
		limits  []*int
		offsets []*int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		empIDs = append(empIDs, query.Inp.(*RoomsInp).EmployeeID)
		limits = append(limits, query.Inp.(*RoomsInp).Params.Limit)
		offsets = append(offsets, query.Inp.(*RoomsInp).Params.Offset)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, orderpos, room_id, name, photo_url, view, last_msg_read, last_msg_id, notify
		FROM load_emp_rooms(
	        $1::text[],
		    $2::bigint[],
		    $3::int[],
		    $4::int[]
		)
	`,
		pq.Array(ptrs),
		pq.Array(empIDs),
		pq.Array(limits),
		pq.Array(offsets),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("rooms:" + err.Desk())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // Каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := new(model.Room)

		if err = rows.Scan(&ptr, &m.Pos, &m.RoomID, &m.Name, &m.PhotoURL, &m.View, &m.LastMessageRead, &m.LastMessageID, &m.Notify); err != nil {
			//c.Dataloader.healer.Alert("rooms (scan rows):" + err.Desk())
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
