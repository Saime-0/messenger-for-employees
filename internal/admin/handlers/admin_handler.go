package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/saime-0/messenger-for-employee/graph/resolver"
	"github.com/saime-0/messenger-for-employee/internal/admin/responder"
	"github.com/saime-0/messenger-for-employee/internal/models"
	"github.com/saime-0/messenger-for-employee/internal/repository"
	"github.com/saime-0/messenger-for-employee/internal/res"
	"log"
	"net/http"
	"time"
)

type AdminHandler struct {
	Resolver *resolver.Resolver
}

func NewAdminHandler(r *mux.Router, resolver *resolver.Resolver) *AdminHandler {
	h := &AdminHandler{Resolver: resolver}
	r.Use()
	h.initAPI(r)
	return h
}

func (h AdminHandler) initAPI(r *mux.Router) {
	adm := r.PathPrefix("/admin").Subrouter()
	h.initEmployeesRoutes(adm)
	h.initRoomsRoutes(adm)
	h.initTagsRoutes(adm)
}

func adminAuth(repo *repository.Repositories) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("started %s %s", r.Method, r.RequestURI)
			token := r.Header.Get(res.AuthHeader)
			if len(token) == 0 {
				responder.Error(w, http.StatusUnauthorized, "missing \"Authorization\" header")
				return
			}
			admin, err := repo.Admins.AdminByToken(token)
			if responder.End(err, w, http.StatusInternalServerError, "bad") {
				return
			}
			r.WithContext(context.WithValue(
				r.Context(),
				res.CtxAdminData,
				admin,
			))
			//a, ok := r.Context().Value(res.CtxAdminData).(*models.Admin)

			next.ServeHTTP(w, r)

		})
	}
}

func Who(r *http.Request) *models.Admin {
	return r.Context().Value(res.CtxAdminData).(*models.Admin)
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responder.Writer{
			ResponseWriter: w,
			Code:           http.StatusOK,
		}
		next.ServeHTTP(rw, r)

		log.Printf(
			"completed with %d %s in %v\n",
			rw.Code,
			http.StatusText(rw.Code),
			time.Since(start),
		)
	})
}
