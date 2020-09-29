// Package handlers has functions to construct HTTP handlers for both the internal and external API.
package handlers

import (
	"bytes"
	"encoding/json"
	"io"
)

func writeJSON(w io.Writer, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		// It's okay to panic here, as everything is clean and marshal'able
		panic(err)
	}
	// Never write null, return an empty array instead
	if bytes.Equal(body, []byte("null")) {
		body = []byte("[]")
	}
	// We ignore this error because there's nothing to do (except perhaps logging)
	_, _ = w.Write(body)
}
