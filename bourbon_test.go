package bourbon

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestBourbonCreateConfig(t *testing.T) {
	c := createConfig()
	r := &router{routes: make(map[string][]Route)}
	assert.Equal(t, "5000", c.Port)
	assert.Equal(t, new(server), c.Server)
	assert.Equal(t, r, c.Router)
}
