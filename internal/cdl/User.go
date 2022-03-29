package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *userResult) isRequestResult() {}
func (r *userInp) isRequestInput()     {}

type (
	userInp struct {
		EmployeeID int
	}
	userResult struct {
		Employee *model.Employee
	}
)

func (d *Dataloader) Employee(employeeID int) (*model.Employee, error) {
	res := <-d.categories.User.addBaseRequest(
		&userInp{
			EmployeeID: employeeID,
		},
		new(userResult),
	)
	if res == nil {
		return nil, d.categories.User.Error
	}
	return res.(*userResult).Employee, nil
}

func (c *parentCategory) user() {
	var (
		inp = c.Requests

		ptrs        []chanPtr
		employeeIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		employeeIDs = append(employeeIDs, query.Inp.(*userInp).EmployeeID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, 
		       coalesce(emp_id, 0), 
		       coalesce(first_name, ''), 
		       coalesce(last_name, ''), 
		       coalesce(joined_at, 0) 
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, employeeid)
		LEFT JOIN employees e ON e.emp_id = inp.employeeid
		`,
		pq.Array(ptrs),
		pq.Array(employeeIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("user:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // Каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := new(model.Employee)
		if err = rows.Scan(&ptr, &m.EmpID, &m.FirstName, &m.LastName, &m.JoinedAt); err != nil {
			//c.Dataloader.healer.Alert("user (scan rows):" + err.Error())
			c.Error = err
			return
		}
		if m.EmpID == 0 {
			m = nil
		}

		request := c.getRequest(ptr)
		request.Result.(*userResult).Employee = m
	}

	c.Error = nil
}
