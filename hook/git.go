package hook

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/ejholmes/hookshot/events"
)

// Checkout checks out the repository and returns the path where it's located.
func Checkout(event events.Push, w io.Writer) (string, error) {
	branch := strings.Replace(event.Ref, "refs/heads/", "", -1)
	sha := event.HeadCommit.ID
	repo := event.Repository.FullName

	dir, err := ioutil.TempDir("", sha)
	if err != nil {
		return dir, err
	}

	cmd := exec.Command("git", "clone", "--depth=50", fmt.Sprintf("--branch=%s", branch), fmt.Sprintf("git://github.com/%s.git", repo), dir)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return dir, err
	}

	cmd = exec.Command("git", "checkout", "-qf", sha)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = dir

	return dir, cmd.Run()
}
