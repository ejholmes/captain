// Pacakge hook provides an http.Handler implementation to handle GitHub push
// webhooks and build the docker images using captain.
package hook

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/ejholmes/hookshot"
	"github.com/ejholmes/hookshot/events"
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

	if err := captain(dir, os.Stdout, "build"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := captain(dir, os.Stdout, "push"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func captain(dir string, w io.Writer, command string) error {
	cmd := exec.Command("captain", command)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = dir
	return cmd.Run()
}
