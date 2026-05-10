package util

import (
	"encoding/json"
	"net/http"

	"github.com/dinno7/ride-sharing/shared/contracts"
)

func NewSuccessPayload(
	data any,
	message string,
) *contracts.APIResponse {
	if len(message) == 0 {
		message = "Process done successfull"
	}
	return &contracts.APIResponse{Ok: true, Data: data, Message: message, Error: nil}
}

func NewErrorPayload(
	code string,
	message string,
) *contracts.APIResponse {
	if len(message) == 0 {
		message = "Process failed"
	}
	return &contracts.APIResponse{
		Data:    nil,
		Ok:      false,
		Message: message,
		Error: &contracts.APIError{
			Code:    code,
			Message: message,
		},
	}
}

func SuccessResponse(w http.ResponseWriter, status int, message string, data any) error {
	return sendJSONResponse(w, status, NewSuccessPayload(data, message))
}

func HttpOkResponse(w http.ResponseWriter, message string, data any) error {
	return SuccessResponse(w, http.StatusOK, message, data)
}

func HttpNoContentResponse(w http.ResponseWriter) error {
	return sendJSONResponse(w, http.StatusNoContent, nil)
}

// --------------------

func ErrorResponse(w http.ResponseWriter, status int, message string) error {
	return sendJSONResponse(w, status, NewErrorPayload(http.StatusText(status), message))
}

func HttpValidationErrorResponse(w http.ResponseWriter, message string) error {
	return ErrorResponse(w, http.StatusBadRequest, message)
}

func sendJSONResponse(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
