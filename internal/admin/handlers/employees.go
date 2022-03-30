package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/admin/request_models"
	"github.com/saime-0/http-cute-chat/internal/admin/responder"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"log"
	"net/http"
)

func (h *AdminHandler) initEmployeesRoutes(r *mux.Router) {
	emp := r.PathPrefix("/emp").Subrouter()
	{
		emp.HandleFunc("/create", h.CreateEmployee).Methods(http.MethodPost)

	}
}

func (h *AdminHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.CreateEmployee{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusBadRequest, "bad")
		return
	}
	inp.Token, _ = utils.HashPassword(inp.Token, h.Resolver.Config.GlobalPasswordSalt)
	log.Printf("%#v", inp) // debug
	id, err := h.Resolver.Services.Repos.Employees.CreateEmployee(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
		return
	}

	responder.Respond(w, http.StatusOK, &request_models.CreateEmployeeResult{EmpID: id})
}
