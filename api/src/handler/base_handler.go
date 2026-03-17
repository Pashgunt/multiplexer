package handler

import (
	"encoding/json"
	"net/http"
	"transport/api/src/dto"
)

type baseHandler struct {
}

func sendJSON(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		return
	}
}

func sendCreated(w http.ResponseWriter, id string) {
	sendJSON(w, http.StatusCreated, dto.ApiResponse{
		Success: true,
		Data: map[string]string{
			"result": id,
		},
	})
}

func sendError(w http.ResponseWriter, status int, errMsg string) {
	sendJSON(w, status, dto.ApiResponse{
		Success: false,
		Error:   errMsg,
	})
}
