package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/admin/request_models"
	"github.com/saime-0/messenger-for-employee/internal/admin/responder"
	"github.com/saime-0/messenger-for-employee/internal/res"
	"log"
	"net/http"
)

func (h *AdminHandler) initRoomsRoutes(r *mux.Router) {
	emp := r.PathPrefix("/rooms").Subrouter()
	{
		emp.HandleFunc("/create", h.CreateRoom).Methods(http.MethodPost)
		emp.HandleFunc("/drop", h.DropRoom).Methods(http.MethodPost)
		emp.HandleFunc("/add-emp", h.AddEmployees).Methods(http.MethodPost)
		emp.HandleFunc("/kick-emp", h.KickEmployees).Methods(http.MethodPost)
	}
}

func (h *AdminHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.CreateRoom{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusBadRequest, "bad")
		return
	}

	id, err := h.Resolver.Services.Repos.Rooms.CreateRoom(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
		return
	}

	responder.Respond(w, http.StatusOK, &request_models.CreateRoomResult{RoomID: id})
}

func (h *AdminHandler) DropRoom(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.DropRoom{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusBadRequest, "bad")
		return
	}
	employees, err := h.Resolver.Services.Repos.Rooms.RoomMembersID(inp.RoomID)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
		return
	}

	err = h.Resolver.Services.Repos.Rooms.DropRoom(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
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

func (h *AdminHandler) AddEmployees(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.AddEmployeeToRooms{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusBadRequest, "bad")
		return
	}
	log.Printf("%#v", inp) // debug

	err = h.Resolver.Services.Repos.Rooms.AddEmployeeToRoom(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
		return
	}
	h.Resolver.Subix.NotifyRoomMembers(
		&model.MemberAction{
			Action:  model.ActionAdd,
			EmpID:   inp.Employee,
			RoomIDs: inp.Rooms,
		},
		inp.Rooms...,
	)
	h.Resolver.Subix.NotifyEmployees(
		&model.MemberAction{
			Action:  model.ActionAdd,
			EmpID:   inp.Employee,
			RoomIDs: inp.Rooms,
		},
		inp.Employee,
	)

	responder.Respond(w, http.StatusOK, res.Success)
}

func (h *AdminHandler) KickEmployees(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.AddOrDeleteEmployeesInRoom{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusBadRequest, "bad")
		return
	}
	err = h.Resolver.Services.Repos.Rooms.KickEmployeesFromRoom(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
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
