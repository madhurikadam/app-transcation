package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madhurikadam/app-transcation/internal/domain"
	"github.com/madhurikadam/app-transcation/pkg/http/controller"
)

type (
	TranscationService interface {
		CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error)
		GetAccount(ctx context.Context, accountID string) (*domain.Account, error)
		CreateTranscation(ctx context.Context, transcation domain.Transcation) (*domain.Transcation, error)
	}

	Gateway struct {
		controller.BaseController
		transcationSvc TranscationService
	}
)

func NewGateway(transcationSvc TranscationService) Gateway {
	return Gateway{
		transcationSvc: transcationSvc,
	}
}

func (g Gateway) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var create domain.AccountReq
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		g.WriteErrorResponseMsg(w, http.StatusBadRequest, "missing or invalid json body")
		return
	}

	account, err := g.transcationSvc.CreateAccount(r.Context(), create.DocumentNumber)
	if err != nil {
		g.WriteErrorResponseMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	g.WriteJSONResponse(w, http.StatusCreated, account)
}

func routeVar(r *http.Request, key string) string {
	if value, ok := mux.Vars(r)[key]; ok {
		return value
	}

	return ""
}

func (g Gateway) GetAccount(w http.ResponseWriter, r *http.Request) {
	account, err := g.transcationSvc.GetAccount(r.Context(), routeVar(r, "id"))
	if err != nil {
		g.WriteErrorResponseMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	g.WriteJSONResponse(w, http.StatusOK, account)
}

func (g Gateway) CreateTranscation(w http.ResponseWriter, r *http.Request) {
	var create domain.Transcation
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		g.WriteErrorResponseMsg(w, http.StatusBadRequest, "missing or invalid json body")
		return
	}

	tx, err := g.transcationSvc.CreateTranscation(r.Context(), create)
	if err != nil {
		g.WriteErrorResponseMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	g.WriteJSONResponse(w, http.StatusOK, tx)
}
