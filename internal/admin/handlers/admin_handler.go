package handlers

import (
	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/graph/resolver"
	"github.com/saime-0/http-cute-chat/internal/admin/responder"
	"log"
	"net/http"
	"time"
)

type AdminHandler struct {
	Resolver *resolver.Resolver
}

func NewAdminHandler(router *mux.Router, resolver *resolver.Resolver) *AdminHandler {
	h := &AdminHandler{Resolver: resolver}
	h.initAPI(router)
	return h
}

func (h AdminHandler) initAPI(r *mux.Router) {
	adm := r.PathPrefix("/admin").Subrouter()
	h.initEmployeesRoutes(adm)
	h.initRoomsRoutes(adm)
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
