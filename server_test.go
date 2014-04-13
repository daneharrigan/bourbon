package bourbon

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestServerRouter(t *testing.T) {
	config = createConfig()
	server := new(server)
	r := &router{routes: make(map[string][]Route)}
	assert.Equal(t, r, server.Router())
}