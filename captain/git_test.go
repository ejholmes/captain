package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"github.com/harbur/captain/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestGitGetRevision(t *testing.T) {
	assert.Equal(t, 7, len(getRevision()), "Git revision should have length 7 chars")
}

func TestGitGetBranch(t *testing.T) {
	assert.Equal(t, "master", getBranch(), "Git branch should be master")
}

func TestGitIsDirty(t *testing.T) {
	assert.Equal(t, false, isDirty(), "Git should not have local changes")
}

func TestGitIsGit(t *testing.T) {
	assert.Equal(t, true, isGit(), "There should be a git repository")
}
