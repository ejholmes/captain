// Pacakge hook provides an http.Handler implementation to handle GitHub push
// webhooks and build the docker images using captain.
package hook

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ejholmes/hookshot"
	"github.com/ejholmes/hookshot/events"
	"github.com/harbur/captain"
)

// NewServer returns a new http.Handler that will handle the `ping` and `push`
// events from GitHub.
func NewServer(secret string) http.Handler {
	r := hookshot.NewRouter()

	auth := func(h http.HandlerFunc) http.Handler {
		return hookshot.Authorize(h, secret)
	}

	// Setup handlers
	r.Handle("ping", auth(Ping))
	r.Handle("push", auth(Push))

	return r
}

// ListenAndServe starts an http server.
func ListenAndServe(addr string, secret string) error {
	return http.ListenAndServe(addr, NewServer(secret))
}

// Ping is an http.HandlerFunc that will handle the `ping` event from GitHub.
func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Push is an http.HandlerFunc that will handle the `push` event from GitHub.
func Push(w http.ResponseWriter, r *http.Request) {
	var event events.Push

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dir, err := Checkout(event, os.Stdout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config := captain.NewConfig("", filepath.Join(dir, "captain.yml"), true)

	if captain.Build(captain.BuildOptions{Config: config}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build
	// Push
}
