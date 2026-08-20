package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespond(t *testing.T) {
	api := &API{}
	w := httptest.NewRecorder()
	api.respond(w, "test data", nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}

func TestRespondError(t *testing.T) {
	api := &API{}
	w := httptest.NewRecorder()
	apiErr := &apiError{
		typ: errorBadParam,
		err: errors.New("bad parameter"),
	}
	api.respondError(w, apiErr, nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", w.Code)
	}
}

func TestRespondSerializationError(t *testing.T) {
	api := &API{}
	w := httptest.NewRecorder()

	// A channel cannot be serialized to JSON, which will trigger an encoding error.
	invalidData := make(chan int)
	api.respond(w, invalidData, nil)

	if w.Code == http.StatusOK {
		t.Errorf("Expected non-200 status code on serialization failure, got %d", w.Code)
	}
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 Internal Server Error, got %d", w.Code)
	}
}
