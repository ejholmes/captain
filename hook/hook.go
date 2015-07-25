// Pacakge hook provides an http.Handler implementation to handle GitHub push
// webhooks and build the docker images using captain.
package hook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/ejholmes/hookshot"
	"github.com/ejholmes/hookshot/events"
)

// Image is the default docker image that will be used to perform the
// build.
var Image = "harbur/captain-builder:latest"

// Write build output to stdout.
var NewLogger logFactory = StdoutLogger

// Logger represents an interface for something that can give us a place to
// stream logs for the build event.
type logFactory func(events.Push) io.Writer

func StdoutLogger(event events.Push) io.Writer {
	return os.Stdout
}

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

	cmd := exec.Command("docker", "run",
		"--privileged=true",
		"--volumes-from=data",
		"-e", fmt.Sprintf("REPOSITORY=%s", event.Repository.FullName),
		"-e", fmt.Sprintf("BRANCH=%s", strings.Replace(event.Ref, "refs/heads/", "", -1)),
		"-e", fmt.Sprintf("SHA=%s", event.HeadCommit.ID),
		Image,
	)

	out := NewLogger(event)
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Run(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
