package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/messenger-for-employee/graph/model"
)

func (r *TagsResult) isRequestResult() {}
func (r *TagsInp) isRequestInput()     {}

type (
	TagsResult struct {
		Tags *model.Tags
	}
	TagsInp struct {
		EmployeeID int
	}
)

func (d *Dataloader) Tags(empID int) (*model.Tags, error) {
	res := <-d.categories.Tags.addBaseRequest(
		&TagsInp{
			EmployeeID: empID,
		},
		&TagsResult{
			Tags: &model.Tags{
				Tags: []*model.Tag{},
			},
		},
	)
	if res == nil {
		return nil, d.categories.Tags.Error
	}
	return res.(*TagsResult).Tags, nil
}

func (c *parentCategory) tags() {
	var (
		inp = c.Requests

		ptrs   []chanPtr
		empIDs []int
	)
	for _, query := range inp {
		empIDs = append(empIDs, query.Inp.(*TagsInp).EmployeeID)
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr,
				coalesce(t.tag_id, 0),
				coalesce(t.name, '')
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, empid)
		LEFT JOIN positions p on p.emp_id = inp.empid
		LEFT JOIN tags t on t.tag_id = p.tag_id
		`,
		pq.Array(ptrs),
		pq.Array(empIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("tags:" + err.Desk())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // Каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := new(model.Tag)

		if err = rows.Scan(&ptr, &m.TagID, &m.Name); err != nil {
			//c.Dataloader.healer.Alert("tags (scan rows):" + err.Desk())
			c.Error = err
			return
		}
		if m.TagID == 0 {
			continue
		}
		request := c.getRequest(ptr)
		request.Result.(*TagsResult).Tags.Tags = append(request.Result.(*TagsResult).Tags.Tags, m)
	}

	c.Error = nil
}
