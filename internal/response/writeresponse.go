package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	StatusCode int `json:"status_code"`
	Data       any `json:"data"`
}

type ErrorResponse struct {
	StatusCode int `json:"status_code"`
	Error      any `json:"error"`
}

// WriteSuccessResponse writes a successful response to w
func WriteSuccessResponse(w http.ResponseWriter, status int, val any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&SuccessResponse{
		StatusCode: status,
		Data:       val,
	})
}

// WriteErrorResponse writes an error to w.
func WriteErrorResponse(w http.ResponseWriter, status int, err error) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&ErrorResponse{
		StatusCode: status,
		Error:      err,
	})
}
