package bourbon

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestRouteCreateMatcher(t *testing.T) {
	r := &route{parent: new(bourbon), pattern: "/test"}
	assert.Equal(t, true, r.regexp == nil)
	assert.Equal(t, true, r.Regexp().MatchString("/test"))
	assert.Equal(t, true, r.regexp != nil)
}

func TestRouteCreateMatcherWithVariables(t *testing.T) {
	r := &route{parent: new(bourbon), pattern: "/test/{id}"}
	assert.Equal(t, true, r.regexp == nil)
	assert.Equal(t, true, r.Regexp().MatchString("/test/1"))
	assert.Equal(t, true, r.Regexp().MatchString("/test/fe1ca51d-f765-4446-a681-77bd02ad07ca"))
	assert.Equal(t, false, r.Regexp().MatchString("/test/2/failed"))
	assert.Equal(t, true, r.regexp != nil)
}

func TestRouteCreateMatcherWithNestedVariables(t *testing.T) {
	r := &route{parent: new(bourbon), pattern: "/one/{one_id}/two/{two_id}"}
	assert.Equal(t, false, r.Regexp().MatchString("/one/1"))
	assert.Equal(t, false, r.Regexp().MatchString("/one/1/two"))
	assert.Equal(t, false, r.Regexp().MatchString("/one/1/two/2/three"))
	assert.Equal(t, true, r.Regexp().MatchString("/one/1/two/2"))
}

func TestRoutePatternWithParents(t *testing.T) {
	parent := new(bourbon)
	parent.SetPrefix("/parent")
	parent.Get("/test", func() {})
	assert.Equal(t, 1, len(parent.Routes()))
	assert.Equal(t, "/parent/test", parent.Routes()[0].Pattern())

	first := new(bourbon)
	first.SetPrefix("/first")
	first.Get("/test", func() {})
	assert.Equal(t, 1, len(first.Routes()))
	assert.Equal(t, "/first/test", first.Routes()[0].Pattern())

	last := new(bourbon)
	last.SetPrefix("/last")
	last.Get("/test", func() {})
	assert.Equal(t, 1, len(last.Routes()))
	assert.Equal(t, "/last/test", last.Routes()[0].Pattern())

	parent.Mount(first)
	first.Mount(last)

	assert.Equal(t, "/parent/test", parent.Routes()[0].Pattern())
	assert.Equal(t, "/parent/first/test", first.Routes()[0].Pattern())
	assert.Equal(t, "/parent/first/last/test", last.Routes()[0].Pattern())
}

func TestRouteMiddlewareWithParents(t *testing.T) {
	parent := new(bourbon)
	parent.Get("/", func() {})
	parent.Use(func() {})

	first := new(bourbon)
	first.Use(func() {})

	last := new(bourbon)
	last.Use(func() {})
	last.Get("/last", func() {})

	parent.Mount(first)
	first.Mount(last)

	r1 := parent.Routes()[0]
	r2 := last.Routes()[0]

	assert.Equal(t, 1, len(r1.Middleware()))
	assert.Equal(t, 3, len(r2.Middleware()))
}
