package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/saime-0/messenger-for-employee/graph/model"
	"github.com/saime-0/messenger-for-employee/internal/admin/request_models"
	"github.com/saime-0/messenger-for-employee/internal/admin/responder"
	"github.com/saime-0/messenger-for-employee/internal/utils"
	"github.com/saime-0/messenger-for-employee/internal/validator"
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
	if responder.End(err, w, http.StatusBadRequest, "bad") {
		return
	}
	inp.Token, _ = utils.HashPassword(inp.Token, h.Resolver.Config.GlobalPasswordSalt)
	fullName := fmt.Sprintf("%s %s", inp.FirstName, inp.LastName)

	if !validator.ValidateEmployeeFullName(fmt.Sprintf("%s %s", inp.FirstName, inp.LastName)) {
		responder.Error(w, http.StatusBadRequest, "invalid employee name")
		return
	}
	emps, err := h.Resolver.Services.Repos.Employees.FindEmployees(&model.FindEmployees{
		EmpID:  nil,
		RoomID: nil,
		TagID:  nil,
		Name:   &fullName,
	})
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}
	if len(emps.Employees) == 0 {
		responder.Error(w, http.StatusBadRequest, "an employee with this name already exists")
		return
	}

	id, err := h.Resolver.Services.Repos.Employees.CreateEmployee(inp)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}

	responder.Respond(w, http.StatusOK, &request_models.CreateEmployeeResult{EmpID: id})
}
