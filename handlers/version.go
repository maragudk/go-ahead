package handlers

import (
	"fmt"
	"io"
	"net/http"
)

// VersionHandler returns the app name and version.
func VersionHandler(name, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, fmt.Sprintf("%v %v", name, version))
	}
}
