package request_models

import "github.com/saime-0/messenger-for-employee/graph/model"

type CreateEmployee struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhotoUrl    string `json:"photo_url"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Token       string `json:"token"`
	Comment     string `json:"comment"`
}
type CreateEmployeeResult struct {
	EmpID int `json:"emp_id"`
}

type CreateRoom struct {
	Name     string         `json:"name"`
	PhotoUrl string         `json:"photo_url"`
	View     model.RoomType `json:"view"`
}
type CreateRoomResult struct {
	RoomID int `json:"room_id"`
}

type DropRoom struct {
	RoomID int `json:"room_id"`
}

type KickEmployeesFromRooms struct {
	Employees []int `json:"emps"`
	Rooms     []int `json:"rooms"`
}

//type AddOrDeleteEmployeesInRoomResult struct {
//	RoomID int `json:"room_id"`
//}

type AddEmployeeToRooms struct {
	EmpID int   `json:"emp_id"`
	Rooms []int `json:"rooms"`
}

type CreateTag struct {
	Name string `json:"name"`
}

type CreateTagResult struct {
	TagID int `json:"tag_id"`
}

type UpdateTag struct {
	TagID int    `json:"tag_id"`
	Name  string `json:"name"`
}

type DropTag struct {
	TagID int `json:"tag_id"`
}

type GiveTag struct {
	EmpID  int   `json:"emp_id"`
	TagIDs []int `json:"tags"`
}

type TakeTag struct {
	EmpID int `json:"emp_id"`
	TagID int `json:"tag_id"`
}
