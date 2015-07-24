package hook

import (
	"net/http/httptest"
	"testing"

	"github.com/ejholmes/hookshot/hooker"
	"github.com/stretchr/testify/assert"
)

// Shared secret to sign GitHub payloads.
const secret = "secret"

func TestPing(t *testing.T) {
	testServer(func(c *hooker.Client) {
		resp, _ := c.Ping(hooker.DefaultPing)
		assert.Equal(t, 200, resp.StatusCode, "Expected a StatusOK response")
	})
}

// testServer starts a new httptest.Server and initializes a hooker client,
// passing it to fn. The test server will be closed when fn returns.
func testServer(fn func(*hooker.Client)) {
	s := httptest.NewServer(NewServer(secret))
	defer s.Close()

	c := hooker.NewClient(nil)
	c.URL = s.URL
	c.Secret = secret

	fn(c)
}
