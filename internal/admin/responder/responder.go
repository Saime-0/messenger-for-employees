package responder

import (
	"encoding/json"
	"log"
	"net/http"
)

type Writer struct {
	http.ResponseWriter
	Code int
}

func (w *Writer) WriteHeader(statusCode int) {
	w.Code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

type RespError struct {
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, code int, msgError string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	if msgError != "" {
		json.NewEncoder(w).Encode(RespError{
			Message: msgError,
		})
	}

}

func End(err error, w http.ResponseWriter, code int, msgError string) (fail bool) {
	if err != nil {
		log.Printf("%#v", err) // debug
		Error(w, code, msgError)
		return true
	}
	return
}
