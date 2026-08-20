package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type API struct{}

type apiError struct {
	typ string
	err error
}

const (
	errorBadParam = "bad_data"
	errorType     = "bad_data"
	errorTimeout  = "timeout"
	errorCanceled = "canceled"
	errorExec     = "execution"
	errorInternal = "internal"
)

type response struct {
	Status    string      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	ErrorType string      `json:"errorType,omitempty"`
	Error     string      `json:"error,omitempty"`
	Warnings  []string    `json:"warnings,omitempty"`
}

func (api *API) respond(w http.ResponseWriter, data interface{}, warnings []string) {
	resp := &response{
		Status:   "success",
		Data:     data,
		Warnings: warnings,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(resp); err != nil {
		api.respondError(w, &apiError{errorInternal, err}, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (api *API) respondError(w http.ResponseWriter, apiErr *apiError, data interface{}) {
	var code int
	switch apiErr.typ {
	case errorBadParam, errorType:
		code = http.StatusBadRequest
	case errorTimeout:
		code = http.StatusServiceUnavailable
	case errorCanceled:
		code = 499
	case errorExec:
		code = http.StatusUnprocessableEntity
	case errorInternal:
		code = http.StatusInternalServerError
	default:
		code = http.StatusInternalServerError
	}

	resp := &response{
		Status:    "error",
		ErrorType: apiErr.typ,
		Error:     apiErr.err.Error(),
		Data:      data,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(resp); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","errorType":"internal","error":"internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(buf.Bytes())
}
