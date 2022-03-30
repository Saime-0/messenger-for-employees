package cdl

import (
	"fmt"
	"github.com/lib/pq"
)

func (r *EmployeeIsRoomMemberResult) isRequestResult() {}
func (r *EmployeeIsRoomMemberInp) isRequestInput()     {}

type (
	EmployeeIsRoomMemberInp struct {
		EmployeeID int
		RoomID     int
	}
	EmployeeIsRoomMemberResult struct {
		Exists bool
	}
)

func (d *Dataloader) EmployeeIsRoomMember(employeeID, roomID int) (bool, error) {

	res := <-d.categories.EmployeeIsRoomMember.addBaseRequest(
		&EmployeeIsRoomMemberInp{
			EmployeeID: employeeID,
			RoomID:     roomID,
		},
		new(EmployeeIsRoomMemberResult),
	)
	if res == nil {
		return false, d.categories.EmployeeIsRoomMember.Error
	}
	return res.(*EmployeeIsRoomMemberResult).Exists, nil
}

func (c *parentCategory) employeeIsRoomMember() {
	var (
		inp = c.Requests

		ptrs        []chanPtr
		employeeIDs []int
		roomIDs     []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		employeeIDs = append(employeeIDs, query.Inp.(*EmployeeIsRoomMemberInp).EmployeeID)
		roomIDs = append(roomIDs, query.Inp.(*EmployeeIsRoomMemberInp).RoomID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, room_id is not null 
		FROM unnest($1::varchar[], $2::bigint[], $3::bigint[]) inp(ptr, empid, roomid)
		LEFT JOIN members m ON m.room_id = inp.roomid AND m.emp_id = inp.empid
		`,
		pq.Array(ptrs),
		pq.Array(employeeIDs),
		pq.Array(roomIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("employeeIsRoomMember:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr      chanPtr
		isMember bool
	)
	for rows.Next() {

		if err = rows.Scan(&ptr, &isMember); err != nil {
			//c.Dataloader.healer.Alert("employeeIsRoomMember (scan rows):" + err.Error())
			c.Error = err
			return
		}

		request := c.getRequest(ptr)
		request.Result.(*EmployeeIsRoomMemberResult).Exists = isMember
	}

	c.Error = nil
}
