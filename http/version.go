package http

import (
	"context"
	"net/http"
)

// MakeVersionHandler makes a HTTP handler that returns information about the version of the component.
func MakeVersionHandler(componentName, ComponentID, version, environment, gitCommit string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		var info = struct {
			Name    string `json:"name"`
			ID      string `json:"id"`
			Version string `json:"version"`
			Env     string `json:"environment"`
			Commit  string `json:"commit"`
		}{
			Name:    componentName,
			ID:      ComponentID,
			Version: version,
			Env:     environment,
			Commit:  gitCommit,
		}
		_ = EncodeReply(context.TODO(), w, info)
	})
}
