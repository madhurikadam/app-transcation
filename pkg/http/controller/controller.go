package controller

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type BaseController struct {
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSONResponse writes the given body as json encoded data and sets the
// Content-Type-Header accordingly
func (b BaseController) WriteJSONResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Error("failed to write response", err)
	}
}

// WriteErrorResponseMsg
func (b BaseController) WriteErrorResponseMsg(w http.ResponseWriter, status int, msg string) {
	b.WriteJSONResponse(w, status, ErrorResponse{
		Error: msg,
	})
}
