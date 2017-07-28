package service

// This file provides server-side bindings for the HTTP transport.
// It utilizes the transport/http.Server.

import (
	"context"
	"crypto/rc4"
	"encoding/json"
	"errors"
	"net/http"
)

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	msg := err.Error()

	switch err {
	case ErrBadRequest:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}

	w.WriteHeader(code)
	s, err := json.Marshal(errorWrapper{Error: msg})
	d := RC4Crypt(s)
	w.Write(d)
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"Err"`
}

func EncodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	s, err := json.Marshal(response)
	d := RC4Crypt(s)
	w.Write(d)
	return err
}

func RC4Crypt(s []byte) []byte {
	key := []byte("e39e7594feaaf4af26cdaea078a316cb")
	c, _ := rc4.NewCipher(key)
	d := make([]byte, len(s))
	c.XORKeyStream(d, s)
	return d
}
