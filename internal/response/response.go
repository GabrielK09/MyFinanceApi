package response

import (
	"encoding/json"
	"net/http"
)

type ResponseData struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    any    `json:"data"`
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func ErrorResponse(message string, data any) ResponseData {
	return ResponseData{
		Status:  false,
		Message: message,
		Data:    data,
	}
}

func SuccessResponse(message string, data any) ResponseData {
	return ResponseData{
		Status:  true,
		Message: message,
		Data:    data,
	}
}
