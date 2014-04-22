package bourbon

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestServerRouter(t *testing.T) {
	server := new(defaultServer)
	r := &defaultRouter{routes: make(map[string][]*Route)}
	assert.Equal(t, r, server.Router())
}
