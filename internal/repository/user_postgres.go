package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"strings"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) CreateUser(userModel *models.RegisterData) (err error) {
	err = r.db.QueryRow(`
		WITH u AS (
			INSERT INTO units (domain, name, type) 
			VALUES ($1, $2, 'USER') 
			RETURNING id
			) 
		INSERT INTO users (id, hashed_password, email) 
		SELECT u.id, $3, $4 
		FROM u 
		RETURNING id`,
		userModel.Domain,
		userModel.Name,
		userModel.HashPassword,
		userModel.Email,
	).Err()

	return
}

// deprecated
func (r *UsersRepo) User(employeeID int) (*model.User, error) {
	user := &model.User{
		Unit: new(model.Unit),
	}
	err := r.db.QueryRow(`
		SELECT id, domain, name, type
		FROM units
		WHERE id = $1`,
		employeeID,
	).Scan(
		&user.Unit.ID,
		&user.Unit.Domain,
		&user.Unit.Name,
		&user.Unit.Type,
	)

	return user, err
}

func (r *UsersRepo) UserExistsByRequisites(inp *models.LoginRequisites) (exists bool, err error) {
	err = r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1 AND hashed_password = $2
		)`,
		inp.Email,
		inp.HashedPasswd,
	).Scan(&exists)

	return

}

func (r *UsersRepo) GetEmployeeIDByRequisites(inp *models.LoginRequisites) (id int, err error) {
	err = r.db.QueryRow(`
		SELECT units.id
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE users.email = $1 AND users.hashed_password = $2`,
		inp.Email,
		inp.HashedPasswd,
	).Scan(&id)

	return
}

func (r *UsersRepo) GetUserByDomain(domain string) (user models.UserInfo, err error) {
	err = r.db.QueryRow(`
		SELECT units.id,units.domain,units.name
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.domain = $1`,
		domain,
	).Scan(
		&user.ID,
		&user.Domain,
		&user.Name,
	)
	if err != nil {
		return // user, err
	}
	return
}

func (r *UsersRepo) GetUserByID(id int) (user models.UserInfo, err error) {
	err = r.db.QueryRow(`
		SELECT units.id,units.domain,units.name
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Domain,
		&user.Name,
	)
	if err != nil {
		return // user, err
	}
	return
}

func (r *UsersRepo) GetCountUserOwnedChats(employeeID int) (count int, err error) {
	err = r.db.QueryRow(`
		SELECT count(*)
		FROM chats 
		WHERE owner_id = $1`,
		employeeID,
	).Scan(&count)
	return
}

func (r *UsersRepo) UserExistsByID(employeeID int) (exists bool) {
	r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE id = $1
		)`,
		employeeID,
	).Scan(&exists)

	return
}

func (r *UsersRepo) UserExistsByDomain(userDomain string) (exists bool) {
	r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM units
			INNER JOIN users
			ON users.id = units.id
			WHERE units.domain = $1
		)`,
		userDomain,
	).Scan(&exists)

	return
}

func (r *UsersRepo) Me(empID int) (*model.Me, error) {
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
		return nil, nil
	}

	return me, err
}

func (r *UsersRepo) OwnedChats(employeeID int) (*model.Chats, error) {
	chats := &model.Chats{
		Chats: []*model.Chat{},
	}
	rows, err := r.db.Query(`
		SELECT units.id, units.domain, units.name, units.type, chats.private
		FROM units INNER JOIN chats 
		ON units.id = chats.id 
		WHERE chats.owner_id = $1`,
		employeeID,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Chat{
			Unit: new(model.Unit),
		}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type, &m.Private); err != nil {
			return nil, err
		}

		chats.Chats = append(chats.Chats, m)
	}

	return chats, nil
}

func (r *UsersRepo) Chats(employeeID int) (*model.Chats, error) {
	chats := &model.Chats{
		Chats: []*model.Chat{},
	}
	rows, err := r.db.Query(`
		SELECT units.id, units.domain, units.name, units.type, chats.private
		FROM units 
		INNER JOIN chats 
			ON units.id = chats.id 
		INNER JOIN chat_members
			ON units.id = chat_members.chat_id
		WHERE chat_members.employee_id = $1`,
		employeeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Chat{
			Unit: new(model.Unit),
		}
		if err = rows.Scan(&m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type, &m.Private); err != nil {
			return nil, err
		}

		chats.Chats = append(chats.Chats, m)
	}

	return chats, nil
}
func (r *UsersRepo) ChatsID(employeeID int) ([]int, error) {
	rows, err := r.db.Query(
		`SELECT chat_id
		FROM chat_members
		WHERE employee_id = $1`,
		employeeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats, err := completeIntArray(rows)

	return chats, err
}
func (r *UsersRepo) FindEmployees(inp *model.FindEmployees) (*model.Employees, error) {
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

func (r UsersRepo) UpdateMe(employeeID int, inp *model.UpdateMeDataInput) (*model.UpdateUser, error) {
	unit := &model.UpdateUser{}
	err := r.db.QueryRow(`
		WITH u AS (
			UPDATE units
			SET 
			    name = COALESCE($2::VARCHAR, name), 
			    domain = COALESCE($3::VARCHAR, domain)
			WHERE id = $1
		    RETURNING domain, name
		)
		UPDATE users
		SET 
		    hashed_password = COALESCE($4::VARCHAR, hashed_password),
		    email = COALESCE($5::VARCHAR, email)
		FROM u
		WHERE id = $1
		RETURNING id, u.domain, u.name
		`,
		employeeID,
		inp.Name,
		inp.Domain,
		inp.Password,
		inp.Email,
	).Scan(
		&unit.ID,
		&unit.Domain,
		&unit.Name,
	)

	return unit, err
}

func (r UsersRepo) GetRegistrationSession(email, code string) (*models.RegisterData, error) {
	regi := &models.RegisterData{}
	err := r.db.QueryRow(`
	    SELECT coalesce(domain, ''), 
	           coalesce(name, ''), 
	           coalesce(email, ''), 
	           coalesce(hashed_password, '')
		FROM (SELECT 1) _
		LEFT JOIN registration_session ON email = $1 AND verify_code = $2
		`,
		email,
		code,
	).Scan(
		&regi.Domain,
		&regi.Name,
		&regi.Email,
		&regi.HashPassword,
	)
	if err != nil {
		return nil, err
	}
	if regi.Domain == "" {
		regi = nil
	}
	return regi, nil
}

func (r UsersRepo) DeleteRegistrationSession(email string) (err error) {
	err = r.db.QueryRow(`
		DELETE FROM registration_session
	    WHERE email = $1
		`,
		email,
	).Err()

	return
}

func (r UsersRepo) DeleteRefreshSession(id int) error {
	err := r.db.QueryRow(`
		DELETE FROM refresh_sessions
	    WHERE id = $1
		`,
		id,
	).Err()

	return err
}

func (r *UsersRepo) CreateRegistrationSession(userModel *models.RegisterData, expAt int64) (verifyCode string, err error) {
	err = r.db.QueryRow(`
		INSERT INTO registration_session (domain, name, email, hashed_password, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING verify_code`,
		userModel.Domain,
		userModel.Name,
		userModel.Email,
		userModel.HashPassword,
		expAt,
	).Scan(&verifyCode)

	return
}

func (r *UsersRepo) EmailIsFree(email string) (free bool, err error) {
	err = r.db.QueryRow(`
		SELECT 
		EXISTS (
			SELECT 1 
			FROM users
			WHERE email = $1
		) 
		OR
		EXISTS (
		    SELECT 1
		    FROM registration_session
		    WHERE email = $1
		)`,
		email,
	).Scan(&free)

	return !free, err
}
