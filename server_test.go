package bourbon

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestServerRouter(t *testing.T) {
	server := new(defaultServer)
	r := createDefaultRouter()
	assert.Equal(t, r, server.Router())
}
