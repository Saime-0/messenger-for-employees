package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/messenger-for-employee/graph/model"
)

func (r *employeeResult) isRequestResult() {}
func (r *employeeInp) isRequestInput()     {}

type (
	employeeInp struct {
		EmployeeID int
	}
	employeeResult struct {
		Employee *model.Employee
	}
)

func (d *Dataloader) Employee(employeeID int) (*model.Employee, error) {
	res := <-d.categories.Employee.addBaseRequest(
		&employeeInp{
			EmployeeID: employeeID,
		},
		new(employeeResult),
	)
	if res == nil {
		return nil, d.categories.Employee.Error
	}
	return res.(*employeeResult).Employee, nil
}

func (c *parentCategory) user() {
	var (
		inp = c.Requests

		ptrs        []chanPtr
		employeeIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		employeeIDs = append(employeeIDs, query.Inp.(*employeeInp).EmployeeID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, 
		       coalesce(id, 0),
		       coalesce(first_name, ''),
		       coalesce(last_name, '')
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, employeeid)
		LEFT JOIN employees e ON e.id = inp.employeeid
		`,
		pq.Array(ptrs),
		pq.Array(employeeIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("user:" + err.Desk())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // Каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := new(model.Employee)
		if err = rows.Scan(&ptr, &m.EmpID, &m.FirstName, &m.LastName); err != nil {
			//c.Dataloader.healer.Alert("user (scan rows):" + err.Desk())
			c.Error = err
			return
		}
		if m.EmpID == 0 {
			m = nil
		}

		request := c.getRequest(ptr)
		request.Result.(*employeeResult).Employee = m
	}

	c.Error = nil
}
