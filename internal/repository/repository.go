package repository

import (
	"database/sql"
)

type Repositories struct {
	Auth      *AuthRepo
	Employees *EmployeesRepo
	Tags      *TagsRepo
	Rooms     *RoomsRepo
	Prepares  *PreparesRepo
	Admins    *AdminsRepo
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Auth:      NewAuthRepo(db),
		Employees: NewEmployeesRepo(db),
		Tags:      NewTagsRepo(db),
		Rooms:     NewRoomsRepo(db),
		Prepares:  NewPreparesRepo(db),
		Admins:    NewAdminsRepo(db),
	}
}
