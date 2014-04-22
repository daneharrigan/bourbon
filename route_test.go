package bourbon

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestRouteCreateMatcher(t *testing.T) {
	r := &Route{Parent: new(bourbon), Pattern: "/test"}
	assert.Equal(t, true, r.Regexp == nil)
	r.Build()
	assert.Equal(t, true, r.Regexp.MatchString("/test"))
}

func TestRouteCreateMatcherWithVariables(t *testing.T) {
	r := &Route{Parent: new(bourbon), Pattern: "/test/{id}"}
	r.Build()
	assert.Equal(t, true, r.Regexp.MatchString("/test/1"))
	assert.Equal(t, true, r.Regexp.MatchString("/test/fe1ca51d-f765-4446-a681-77bd02ad07ca"))
	assert.Equal(t, false, r.Regexp.MatchString("/test/2/failed"))
}

func TestRouteCreateMatcherWithNestedVariables(t *testing.T) {
	r := &Route{Parent: new(bourbon), Pattern: "/one/{one_id}/two/{two_id}"}
	r.Build()
	assert.Equal(t, false, r.Regexp.MatchString("/one/1"))
	assert.Equal(t, false, r.Regexp.MatchString("/one/1/two"))
	assert.Equal(t, false, r.Regexp.MatchString("/one/1/two/2/three"))
	assert.Equal(t, true, r.Regexp.MatchString("/one/1/two/2"))
}
