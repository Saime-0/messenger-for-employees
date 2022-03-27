package models

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/pkg/scheduler"
)

type Invite struct {
	Code   string `json:"code"`
	ChatID int    `json:"chat_id,omitempty"`
	Aliens int    `json:"aliens"`
	Exp    int64  `json:"exp"`
}

type CreateMessage struct {
	TargetMsgID *int
	EmployeeID  int
	RoomID      int
	Body        string
	// CreatedAt int64 migrate to postgres
}

/* // todo
msg_type:
	- system message
		- sender
		- event
		- body
	- user message
	- formatted message
msg_fields:
	- text
	- photo
	- file
	- vote
	- music
	- video
*/

type RoleReference struct {
	ID    *int
	Name  *string
	Color *model.HexColor
}

type RefreshSession struct {
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	ExpAt        int64
}

type Unit struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Type   res.UnitType
}

type ScheduleInvite struct {
	ChatID int
	Code   string
	Exp    *int64
	Task   *scheduler.Task
}

type ScheduleRegisterSession struct {
	Email string
	Exp   int64
	Task  *scheduler.Task
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

type SubUser struct {
	MemberID *int
	ChatID   *int
}
