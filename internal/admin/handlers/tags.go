package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/saime-0/messenger-for-employee/internal/admin/request_models"
	"github.com/saime-0/messenger-for-employee/internal/admin/responder"
	"github.com/saime-0/messenger-for-employee/internal/res"
	"log"
	"net/http"
)

func (h *AdminHandler) initTagsRoutes(r *mux.Router) {
	emp := r.PathPrefix("/tags").Subrouter()
	{
		emp.HandleFunc("/create", h.CreateTag).Methods(http.MethodPost)
		emp.HandleFunc("/drop", h.DropTag).Methods(http.MethodPost)
		emp.HandleFunc("/update", h.UpdateTag).Methods(http.MethodPost)
		emp.HandleFunc("/give", h.GiveTag).Methods(http.MethodPost)
		emp.HandleFunc("/take", h.TakeTag).Methods(http.MethodPost)
	}
}

func (h *AdminHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.CreateTag{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusBadRequest, "bad")
		return
	}

	id, err := h.Resolver.Services.Repos.Tags.CreateTag(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
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

	err = h.Resolver.Services.Repos.Tags.DropTag(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
		return
	}

	// todo NotifyAllEmployees
	//h.Resolver.Subix.NotifyEmployees(
	//	&model.DropRoom{
	//		RoomID: inp.RoomID,
	//	},
	//	employees...,
	//)

	responder.Respond(w, http.StatusOK, res.Success)
}

func (h *AdminHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.UpdateTag{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusBadRequest, "bad")
		return
	}

	err = h.Resolver.Services.Repos.Tags.UpdateTag(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
		return
	}

	responder.Respond(w, http.StatusOK, res.Success)
}

func (h *AdminHandler) GiveTag(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.GiveTag{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusBadRequest, "bad")
		return
	}
	log.Printf("%#v", inp) // debug

	err = h.Resolver.Services.Repos.Tags.GiveTag(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
		return
	}

	responder.Respond(w, http.StatusOK, res.Success)
}

func (h *AdminHandler) TakeTag(w http.ResponseWriter, r *http.Request) {
	inp := &request_models.TakeTag{}
	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusBadRequest, "bad")
		return
	}
	log.Printf("%#v", inp) // debug

	err = h.Resolver.Services.Repos.Tags.TakeTag(inp)
	if err != nil {
		log.Println(err) // debug
		responder.Error(w, http.StatusInternalServerError, "bad")
		return
	}

	responder.Respond(w, http.StatusOK, res.Success)
}
