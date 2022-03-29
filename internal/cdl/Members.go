package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *MembersResult) isRequestResult() {}
func (r *MembersInp) isRequestInput()     {}

type (
	MembersResult struct {
		Members *model.Members
	}
	MembersInp struct {
		RoomID int
	}
)

func (d *Dataloader) Members(roomID int) (*model.Members, error) {
	res := <-d.categories.Members.addBaseRequest(
		&MembersInp{
			RoomID: roomID,
		},
		&MembersResult{
			Members: &model.Members{
				Members: []*model.Member{},
			},
		},
	)
	if res == nil {
		return nil, d.categories.Members.Error
	}
	return res.(*MembersResult).Members, nil
}

func (c *parentCategory) members() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		roomIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		roomIDs = append(roomIDs, query.Inp.(*MembersInp).RoomID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr,
				coalesce(m.emp_id, 0),
				coalesce(m.room_id, 0)
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, roomid)
		LEFT JOIN members m ON m.room_id = inp.roomid
		`,
		pq.Array(ptrs),
		pq.Array(roomIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("members:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // Каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := &model.Member{
			Employee: new(model.Employee),
			Room:     new(model.Room),
		}

		if err = rows.Scan(&ptr, &m.Employee.EmpID, &m.Room.RoomID); err != nil {
			//c.Dataloader.healer.Alert("members (scan rows):" + err.Error())
			c.Error = err
			return
		}
		if m.Employee.EmpID == 0 {
			continue
		}
		request := c.getRequest(ptr)
		request.Result.(*MembersResult).Members.Members = append(request.Result.(*MembersResult).Members.Members, m)
	}

	c.Error = nil
}
