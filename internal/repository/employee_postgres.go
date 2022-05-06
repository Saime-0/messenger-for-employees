package repository

import (
	"database/sql"
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/admin/request_models"
	"github.com/saime-0/messenger-for-employee/internal/cerrors"
	"github.com/saime-0/messenger-for-employee/internal/models"
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
		SELECT coalesce(e.id, 0),
		       coalesce(e.first_name,''),
		       coalesce(e.last_name, ''),
		       coalesce(e.email, ''),
		       coalesce(e.phone_number, '')
		FROM employees e
		WHERE e.id = $1`,
		empID,
	).Scan(
		&me.Employee.EmpID,
		&me.Employee.FirstName,
		&me.Employee.LastName,
		&me.Personal.Email,
		&me.Personal.PhoneNumber,
	)
	if me.Employee.EmpID == 0 {
		return nil, cerrors.New("user not found")
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
		SELECT e.id, e.first_name, e.last_name
		FROM employees e 
		    LEFT JOIN positions p ON e.id = p.emp_id 
			LEFT JOIN members m ON m.emp_id = e.id 
		WHERE (
		    $1::BIGINT IS NULL OR
		    e.id = $1::BIGINT
		) AND (
			$2::BIGINT IS NULL OR
			m.room_id = $2::BIGINT
		)
		AND (
			$3::BIGINT IS NULL OR
			p.tag_id = $3::BIGINT
		) AND (
		    $4::VARCHAR = ' ' OR 
		    e.first_name ILIKE $4::VARCHAR OR e.last_name ILIKE $4::VARCHAR
		) AND (
		    $5::VARCHAR = ' ' OR 
		    e.first_name ILIKE $5::VARCHAR OR e.last_name ILIKE $5::VARCHAR
		)
		GROUP BY e.id 
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
		if err = rows.Scan(&m.EmpID, &m.FirstName, &m.LastName); err != nil {
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
			WHERE email = $1 AND password_hash = $2
		)`,
		inp.Email,
		inp.HashedPasswd,
	).Scan(&exists)

	return

}

// returning id=0 if not found
func (r *EmployeesRepo) GetEmployeeIDByRequisites(inp *models.LoginRequisites) (id int, err error) {
	err = r.db.QueryRow(`
		SELECT coalesce((
			SELECT id
			FROM employees
			WHERE email = $1 AND password_hash = $2
		), 0)`,
		inp.Email,
		inp.HashedPasswd,
	).Scan(&id)

	return
}

func (r EmployeesRepo) CreateEmployee(emp *request_models.CreateEmployee) (empID int, err error) {
	err = r.db.QueryRow(`
		INSERT INTO employees (first_name, last_name, email, phone_number, password_hash, comment) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
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
