package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"log/slog"
)

// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

type HttpError struct {
	Status int
	Text   string
}

func (mr *HttpError) Error() string {
	return mr.Text
}

type JsonBodyDecoder struct {
	MaxBytes              int64
	DisallowUnknownFields bool
}

func (d *JsonBodyDecoder) DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &HttpError{Status: http.StatusUnsupportedMediaType, Text: msg}
		}
	}

	if d.MaxBytes > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, d.MaxBytes)
	}

	dec := json.NewDecoder(r.Body)

	if d.DisallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &HttpError{Status: http.StatusBadRequest, Text: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &HttpError{Status: http.StatusBadRequest, Text: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &HttpError{Status: http.StatusBadRequest, Text: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &HttpError{Status: http.StatusBadRequest, Text: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &HttpError{Status: http.StatusBadRequest, Text: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &HttpError{Status: http.StatusRequestEntityTooLarge, Text: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &HttpError{Status: http.StatusBadRequest, Text: msg}
	}

	return nil
}

func JSONBody[T any](
	log *slog.Logger,
	decoder *JsonBodyDecoder,
	w http.ResponseWriter,
	r *http.Request,
) (T, *HttpError) {
	var dst T
	if err := decoder.DecodeJSONBody(w, r, &dst); err != nil {
		var mr *HttpError
		if errors.As(err, &mr) {
			return dst, mr
		}
		log.LogAttrs(r.Context(), slog.LevelError, "failed to decode request body", slog.String("error", err.Error()))
		mr.Status = http.StatusInternalServerError
		mr.Text = http.StatusText(http.StatusInternalServerError)
		return dst, mr
	}
	return dst, nil
}
