package models

import (
	"github.com/saime-0/messenger-for-employee/pkg/scheduler"
)

type CreateMessage struct {
	TargetMsgID *int
	EmployeeID  int
	RoomID      int
	Body        string
	// CreatedAt int64 migrate to postgres
}

type RefreshSession struct {
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	ExpAt        int64
}

type ScheduleRefreshSession struct {
	ID   int
	Exp  int64
	Task *scheduler.Task
}

type LoginRequisites struct {
	Email        string
	HashedPasswd string
}

type RegisterData struct {
	Domain       string
	Name         string
	Email        string
	HashPassword string
}

type SubRoom struct {
	RoomID   *int
	IsMember bool
}
type RoomsAndEmployees struct {
	RoomIDs []int
	EmpIDs  []int
}
