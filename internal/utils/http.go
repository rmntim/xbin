package utils

import (
	"encoding/json"
	"net/http"
)

type httpError struct {
	Message string `json:"message"`
}

func RespondJSON[T any](w http.ResponseWriter, status int, value T) error {
	enc := json.NewEncoder(w)
	w.WriteHeader(status)

	return enc.Encode(value)
}

func MustRespondJSON[T any](w http.ResponseWriter, status int, value T) {
	if err := RespondJSON(w, status, value); err != nil {
		panic(err)
	}
}

func ReadJSON[T any](r *http.Request) (T, error) {
	dec := json.NewDecoder(r.Body)

	var value T

	if err := dec.Decode(&value); err != nil {
		return value, err
	}

	return value, nil
}

func RespondError(w http.ResponseWriter, status int, msg string) error {
	err := httpError{
		Message: msg,
	}
	return RespondJSON(w, status, err)
}

func MustRespondError(w http.ResponseWriter, status int, msg string) {
	if err := RespondError(w, status, msg); err != nil {
		panic(err)
	}
}
