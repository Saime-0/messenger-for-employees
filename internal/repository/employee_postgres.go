package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/admin/request_models"
	"github.com/saime-0/http-cute-chat/internal/models"
	"strings"
)

type EmployeesRepo struct {
	db *sql.DB
}

func NewEmployeesRepo(db *sql.DB) *EmployeesRepo {
	return &EmployeesRepo{
		db: db,
	}
}

func (r *EmployeesRepo) Me(empID int) (*model.Me, error) {
	me := &model.Me{
		Employee: new(model.Employee),
		Personal: new(model.PersonalData),
	}
	err := r.db.QueryRow(`
		SELECT coalesce(e.emp_id, 0), 
		       coalesce(e.first_name,''), 
		       coalesce(e.last_name, ''), 
		       coalesce(e.joined_at, 0),
		       coalesce(e.email, ''), 
		       coalesce(e.phone_number, ''), 
		       coalesce(e.token, '')
		FROM employees e
		WHERE e.emp_id = $1`,
		empID,
	).Scan(
		&me.Employee.EmpID,
		&me.Employee.FirstName,
		&me.Employee.LastName,
		&me.Employee.JoinedAt,
		&me.Personal.Email,
		&me.Personal.PhoneNumber,
		&me.Personal.Token,
	)
	if me.Employee.EmpID == 0 {
		return nil, err
	}

	return me, err
}

func (r *EmployeesRepo) FindEmployees(inp *model.FindEmployees) (*model.Employees, error) {
	var (
		users        = new(model.Employees)
		fullName     []string
		fname, lname = " ", " "
	)
	if inp.Name != nil {
		//*inp.Name = "%" + *inp.Name + "%"
		fullName = strings.Split(*inp.Name, " ")
		fname = "%" + fullName[0] + "%"
		if len(fullName) > 1 {
			lname = "%" + fullName[1] + "%"
		}
	}
	rows, err := r.db.Query(`
		SELECT e.emp_id, e.first_name, e.last_name, e.joined_at
		FROM employees e 
		    LEFT JOIN positions p ON e.emp_id = p.emp_id 
			LEFT JOIN members m ON m.emp_id = e.emp_id 
		WHERE (
		    $1 IS NULL OR
		    e.emp_id = $1
		) AND (
			$2 IS NULL OR
			m.room_id = $2
		)
		AND (
			$3 IS NULL OR
			p.tag_id = $3
		) AND (
		    $4 = ' ' OR 
		    e.first_name ILIKE $4 OR e.last_name ILIKE $4
		) AND (
		    $5 = ' ' OR 
		    e.first_name ILIKE $5 OR e.last_name ILIKE $5
		)
		GROUP BY e.emp_id 
		`,
		inp.EmpID,
		inp.RoomID,
		inp.TagID,
		fname,
		lname,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := new(model.Employee)
		if err = rows.Scan(&m.EmpID, &m.FirstName, &m.LastName, &m.JoinedAt); err != nil {
			return nil, err
		}
		users.Employees = append(users.Employees, m)
	}

	return users, nil
}

func (r EmployeesRepo) DeleteRefreshSession(id int) error {
	err := r.db.QueryRow(`
		DELETE FROM refresh_sessions
	    WHERE id = $1
		`,
		id,
	).Err()

	return err
}

func (r *EmployeesRepo) EmailIsFree(email string) (free bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS (
		    SELECT 1
		    FROM employees
		    WHERE email = $1
		)`,
		email,
	).Scan(&free)

	return !free, err
}

func (r *EmployeesRepo) EmployeeExistsByRequisites(inp *models.LoginRequisites) (exists bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM employees
			WHERE email = $1 AND token = $2
		)`,
		inp.Email,
		inp.HashedPasswd,
	).Scan(&exists)

	return

}

func (r *EmployeesRepo) GetEmployeeIDByRequisites(inp *models.LoginRequisites) (id int, err error) {
	err = r.db.QueryRow(`
		SELECT emp_id
		FROM employees
		WHERE email = $1 AND token = $2`,
		inp.Email,
		inp.HashedPasswd,
	).Scan(&id)

	return
}

func (r EmployeesRepo) CreateEmployee(emp *request_models.CreateEmployee) (empID int, err error) {
	err = r.db.QueryRow(`
		INSERT INTO employees (first_name, last_name, email, phone_number, token, comment) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING emp_id
	`,
		emp.FirstName,
		emp.LastName,
		emp.Email,
		emp.PhoneNumber,
		emp.Token,
		emp.Comment,
	).Scan(&empID)
	return
}