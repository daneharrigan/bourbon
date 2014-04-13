package bourbon

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestRouteCreateMatcher(t *testing.T) {
	r := &route{parent: new(bourbon), pattern: "/test"}
	assert.Equal(t, true, r.matcher == nil)
	assert.Equal(t, true, r.MatchString("/test"))
	assert.Equal(t, true, r.matcher != nil)
}

func TestRouteCreateMatcherWithVariables(t *testing.T) {
	r := &route{parent: new(bourbon), pattern: "/test/{id}"}
	assert.Equal(t, true, r.matcher == nil)
	assert.Equal(t, true, r.MatchString("/test/1"))
	assert.Equal(t, true, r.MatchString("/test/fe1ca51d-f765-4446-a681-77bd02ad07ca"))
	assert.Equal(t, false, r.MatchString("/test/2/failed"))
	assert.Equal(t, true, r.matcher != nil)
}

func TestRouteCreateMatcherWithNestedVariables(t *testing.T) {
	r := &route{parent: new(bourbon), pattern: "/one/{one_id}/two/{two_id}"}
	assert.Equal(t, false, r.MatchString("/one/1"))
	assert.Equal(t, false, r.MatchString("/one/1/two"))
	assert.Equal(t, false, r.MatchString("/one/1/two/2/three"))
	assert.Equal(t, true, r.MatchString("/one/1/two/2"))
}
