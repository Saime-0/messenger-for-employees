package handlers

import (
	"encoding/json"
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

func (h *AdminHandler) initTagsRoutes(r *mux.Router) {
	emp := r.PathPrefix("/tags").Subrouter()
	{ // todo
		emp.HandleFunc("/create", h.CreateTag).Methods(http.MethodPost) // /tags/create
		emp.HandleFunc("/drop", h.DropTag).Methods(http.MethodPost)     // /tags/{tag-id}/drop
		emp.HandleFunc("/update", h.UpdateTag).Methods(http.MethodPost) // /tags/{tag-id}/update
		emp.HandleFunc("/give", h.GiveTag).Methods(http.MethodPost)     // /tags/{tag-id}/give/{emp-id}
		emp.HandleFunc("/take", h.TakeTag).Methods(http.MethodPost)     // /tags/{tag-id}/take/{emp-id}
	}
}

func (h *AdminHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.CreateTag{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if responder.End(err, w, http.StatusBadRequest, "bad") {
		return
	}
	inp.Name = kit.UnitWhitespaces(inp.Name)
	if !validator.ValidateTagName(inp.Name) {
		responder.Error(w, http.StatusBadRequest, "invalid tag name")
		return
	}

	exists, err := h.Resolver.Services.Repos.Tags.TagExistsByName(inp.Name)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}
	if exists {
		responder.Error(w, http.StatusBadRequest, "tag with this name already exists")
		return
	}

	id, err := h.Resolver.Services.Repos.Tags.CreateTag(inp)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}

	responder.Respond(w, http.StatusOK, &request_models.CreateTagResult{TagID: id})
}

func (h *AdminHandler) DropTag(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.DropTag{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusBadRequest, "bad")
		return
	}

	exists, err := h.Resolver.Services.Repos.Tags.TagExistsByID(inp.TagID)
	if !exists {
		responder.Error(w, http.StatusBadRequest, "tag is not exists")
		return
	}

	err = h.Resolver.Services.Repos.Tags.DropTag(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
		return
	}

	h.Resolver.Subix.NotifyAllEmployees(
		&model.DropTag{
			TagID: inp.TagID,
		},
	)

	responder.Respond(w, http.StatusOK, res.Success)
}

func (h *AdminHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.UpdateTag{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if responder.End(err, w, http.StatusBadRequest, "bad") {
		return
	}
	inp.Name = kit.UnitWhitespaces(inp.Name)
	if !validator.ValidateTagName(inp.Name) {
		responder.Error(w, http.StatusBadRequest, "invalid tag name")
		return
	}
	exists, err := h.Resolver.Services.Repos.Tags.TagExistsByID(inp.TagID)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}
	if !exists {
		responder.Error(w, http.StatusBadRequest, "tag is not exists")
		return
	}

	err = h.Resolver.Services.Repos.Tags.UpdateTag(inp)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}

	responder.Respond(w, http.StatusOK, res.Success)
}

func (h *AdminHandler) GiveTag(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.GiveTag{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if responder.End(err, w, http.StatusBadRequest, "bad") {
		return
	}
	log.Printf("%#v", inp) // debug
	if len(inp.TagIDs) == 0 {
		responder.Error(w, http.StatusBadRequest, "\"tags\" must contain tag IDs")
		return
	}
	emps, err := h.Resolver.Services.Repos.Employees.FindEmployees(&model.FindEmployees{
		EmpID: &inp.EmpID,
	})
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}
	if len(emps.Employees) == 0 {
		responder.Error(w, http.StatusBadRequest, "an employee with this name already exists")
		return
	}

	for _, tag := range inp.TagIDs {
		exists, err := h.Resolver.Services.Repos.Tags.TagExistsByID(tag)
		if responder.End(err, w, http.StatusInternalServerError, "bad") {
			return
		}
		if !exists {
			responder.Error(w, http.StatusBadRequest, "tag is not exists")
			return
		}
	}
	err = h.Resolver.Services.Repos.Tags.GiveTag(inp)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}

	responder.Respond(w, http.StatusOK, res.Success)
}

func (h *AdminHandler) TakeTag(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.TakeTag{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if responder.End(err, w, http.StatusBadRequest, "bad") {
		return
	}
	log.Printf("%#v", inp) // debug

	err = h.Resolver.Services.Repos.Tags.TakeTag(inp)
	if responder.End(err, w, http.StatusInternalServerError, "bad") {
		return
	}

	responder.Respond(w, http.StatusOK, res.Success)
}
