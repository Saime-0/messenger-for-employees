package request_models

import "github.com/99designs/gqlgen/integration/models-go"

type CreateEmployee struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Token       string `json:"token"`
	Comment     string `json:"comment"`
}
type CreateEmployeeResult struct {
	EmpID int `json:"emp_id"`
}

type CreateRoom struct {
	Name int           `json:"name"`
	View models.Viewer `json:"view"`
}
type CreateRoomResult struct {
	RoomID int `json:"room_id"`
}

type AddOrDeleteEmployeesInRoom struct {
	Rooms     []int `json:"rooms"`
	Employees []int `json:"employees"`
}

//type AddOrDeleteEmployeesInRoomResult struct {
//	RoomID int `json:"room_id"`
//}

type AddEmployeeToRooms struct {
	Employee int   `json:"emp"`
	Rooms    []int `json:"rooms"`
}
