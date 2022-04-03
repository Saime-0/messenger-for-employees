package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/admin/request_models"
	"github.com/saime-0/messenger-for-employee/internal/admin/responder"
	"github.com/saime-0/messenger-for-employee/internal/res"
	"github.com/saime-0/messenger-for-employee/internal/validator"
	"github.com/saime-0/messenger-for-employee/pkg/kit"
	"log"
	"net/http"
)

func (h *AdminHandler) initRoomsRoutes(r *mux.Router) {
	emp := r.PathPrefix("/rooms").Subrouter()
	{ // todo
		emp.HandleFunc("/create", h.CreateRoom).Methods(http.MethodPost)  //         /rooms/create
		emp.HandleFunc("/drop", h.DropRoom).Methods(http.MethodPost)      //             /rooms/{room-id}/drop
		emp.HandleFunc("/add", h.AddEmployee).Methods(http.MethodPost)    //      [...]
		emp.HandleFunc("/kick", h.KickEmployees).Methods(http.MethodPost) //    [...]
	}
}

func (h *AdminHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.CreateRoom{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if responder.End(err, w, http.StatusBadRequest, "bad") {
		return
	}
	inp.Name = kit.UnitWhitespaces(inp.Name)

	if !inp.View.IsValid() {
		responder.Error(w, http.StatusBadRequest, fmt.Sprintf("available views: %s", model.AllRoomType))
		return
	}
	if !validator.ValidateRoomName(inp.Name) {
		responder.Error(w, http.StatusBadRequest, "invalid room name")
		return
	}
	id, err := h.Resolver.Services.Repos.Rooms.CreateRoom(inp)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}

	responder.Respond(w, http.StatusOK, &request_models.CreateRoomResult{RoomID: id})
}

func (h *AdminHandler) DropRoom(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.DropRoom{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if responder.End(err, w, http.StatusBadRequest, "bad") {
		return
	}

	exists, err := h.Resolver.Services.Repos.Rooms.RoomExists(inp.RoomID)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}
	if !exists {
		responder.Error(w, http.StatusBadRequest, "room is not exists")
		return
	}

	employees, err := h.Resolver.Services.Repos.Rooms.RoomMembersID(inp.RoomID)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}

	err = h.Resolver.Services.Repos.Rooms.DropRoom(inp)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}

	h.Resolver.Subix.NotifyEmployees(
		&model.DropRoom{
			RoomID: inp.RoomID,
		},
		employees...,
	)

	responder.Respond(w, http.StatusOK, res.Success)
}

func (h *AdminHandler) AddEmployee(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.AddEmployeeToRooms{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if responder.End(err, w, http.StatusBadRequest, "bad") {
		return
	}
	log.Printf("%#v", inp) // debug
	emps, err := h.Resolver.Services.Repos.Employees.FindEmployees(&model.FindEmployees{
		EmpID: &inp.EmpID,
	})
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}
	if len(emps.Employees) == 0 {
		responder.Error(w, http.StatusBadRequest, fmt.Sprintf("employee(id:%d) does not exist", inp.EmpID))
		return
	}
	for _, room := range inp.Rooms {
		exists, err := h.Resolver.Services.Repos.Rooms.RoomExists(room)
		if responder.End(err, w, http.StatusInternalServerError, "bad") {
			return
		}
		if !exists {
			responder.Error(w, http.StatusBadRequest, "room is not exists")
			return
		}
	}

	err = h.Resolver.Services.Repos.Rooms.AddEmployeeToRoom(inp)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}
	h.Resolver.Subix.NotifyRoomMembers(
		&model.MemberAction{
			Action:  model.ActionAdd,
			EmpID:   inp.EmpID,
			RoomIDs: inp.Rooms,
		},
		inp.Rooms...,
	)
	h.Resolver.Subix.NotifyEmployees(
		&model.MemberAction{
			Action:  model.ActionAdd,
			EmpID:   inp.EmpID,
			RoomIDs: inp.Rooms,
		},
		inp.EmpID,
	)

	responder.Respond(w, http.StatusOK, res.Success)
}

func (h *AdminHandler) KickEmployees(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.KickEmployeesFromRooms{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if responder.End(err, w, http.StatusBadRequest, "bad") {
		return
	}
	if len(inp.Employees) == 0 {
		responder.Error(w, http.StatusBadRequest, "\"emps\" must contain employee IDs")
		return
	}
	if len(inp.Rooms) == 0 {
		responder.Error(w, http.StatusBadRequest, "\"rooms\" must contain room IDs")
		return
	}
	for _, empID := range inp.Employees {
		emps, err := h.Resolver.Services.Repos.Employees.FindEmployees(&model.FindEmployees{
			EmpID: &empID,
		})
		if responder.End(err, w, http.StatusInternalServerError, "bad") {
			return
		}
		if len(emps.Employees) == 0 {
			responder.Error(w, http.StatusBadRequest,
				fmt.Sprintf("employee(id:%d) is not exists", empID))
			return
		}
	}
	err = h.Resolver.Services.Repos.Rooms.KickEmployeesFromRoom(inp)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}
	for _, empID := range inp.Employees {
		h.Resolver.Subix.NotifyRoomMembers(
			&model.MemberAction{
				Action:  model.ActionDel,
				EmpID:   empID,
				RoomIDs: inp.Rooms,
			},
			inp.Rooms...,
		)
	}

	responder.Respond(w, http.StatusOK, res.Success)
}
