package handler

import (
	"encoding/json"
	"net/http"
	"transport/api/src/dto"
)

type baseHandler struct {
}

func (h *baseHandler) sendJSON(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		return
	}
}

func (h *baseHandler) sendCreated(w http.ResponseWriter, id string) {
	h.sendJSON(w, http.StatusCreated, dto.APIResponse{
		Success: true,
		Data: map[string]string{
			"result": id,
		},
	})
}

func (h *baseHandler) sendItem(w http.ResponseWriter, item interface{}) {
	h.sendJSON(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"result": item,
		},
	})
}

func (h *baseHandler) sendDeleted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *baseHandler) sendError(w http.ResponseWriter, status int, errMsg string) {
	h.sendJSON(w, status, dto.APIResponse{
		Success: false,
		Error:   errMsg,
	})
}
