package repository

import (
	"database/sql"
)

type Repositories struct {
	Auth      *AuthRepo
	Employees *EmployeesRepo
	Tags      *TagsRepo
	Rooms     *RoomsRepo
	Messages  *MessagesRepo
	Prepares  *PreparesRepo
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Auth:      NewAuthRepo(db),
		Employees: NewEmployeesRepo(db),
		Tags:      NewTagsRepo(db),
		Rooms:     NewRoomsRepo(db),
		Messages:  NewMessagesRepo(db),
		Prepares:  NewPreparesRepo(db),
	}
}
