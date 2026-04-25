package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"transport/api/src/command"
	apidto "transport/api/src/dto"
	"transport/api/src/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TargetServiceHandler struct {
	service service.ITargetServiceService
}

func NewTargetServiceHandler(service service.ITargetServiceService) *TargetServiceHandler {
	return &TargetServiceHandler{service: service}
}

func (h TargetServiceHandler) Create(c *gin.Context) {
	var dto apidto.TargetServiceDto

	if err := json.NewDecoder(c.Request.Body).Decode(&dto); err != nil {
		c.JSON(http.StatusBadRequest, apidto.APIResponse{Error: err.Error()})

		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := h.service.Create(ctx, command.CreateTargetServiceCommand{Dto: dto})

	if err != nil {
		c.JSON(http.StatusBadRequest, apidto.APIResponse{Error: err.Error()})

		return
	}

	c.JSON(http.StatusCreated, apidto.APIResponse{Data: id})
}

func (h TargetServiceHandler) Delete(c *gin.Context) {
	err := h.service.Delete(
		c.Request.Context(),
		command.DeleteTargetServiceCommand{ID: uuid.MustParse(c.Param("id"))},
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, apidto.APIResponse{Error: err.Error()})

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h TargetServiceHandler) Get(c *gin.Context) {
	result, err := h.service.Get(
		c.Request.Context(),
		command.GetTargetServiceCommand{ID: uuid.MustParse(c.Param("id"))},
	)

	if result == nil {
		c.JSON(http.StatusNotFound, apidto.APIResponse{Error: "Not Found"})

		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, apidto.APIResponse{Error: err.Error()})

		return
	}

	c.JSON(http.StatusOK, apidto.APIResponse{Data: result})
}
