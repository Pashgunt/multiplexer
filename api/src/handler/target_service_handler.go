package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"transport/api/src/command"
	apidto "transport/api/src/dto"
	"transport/api/src/service"
)

type TargetServiceHandler struct {
	service service.ITargetServiceService
	baseHandler
}

func NewTargetServiceHandler(service service.ITargetServiceService) *TargetServiceHandler {
	return &TargetServiceHandler{service: service}
}

func (h TargetServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto apidto.TargetServiceDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	id, err := h.service.Create(ctx, command.CreateTargetServiceCommand{Dto: dto})

	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())

		return
	}

	sendCreated(w, id)
}
