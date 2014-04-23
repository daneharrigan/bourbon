package bourbon

import "regexp"

var variable = regexp.MustCompile(`{([^}]+)}`)

type route struct {
	parent  Bourbon
	method  string
	pattern string
	handler Handler
	regexp  *regexp.Regexp
}

func (r *route) SetParent(b Bourbon) {
	r.parent = b
}

func (r *route) Parent() Bourbon {
	return r.parent
}

func (r *route) Method() string {
	return r.method
}

func (r *route) Pattern() string {
	pattern := r.pattern
	parent := r.Parent()

	for parent != nil {
		pattern = parent.Prefix() + pattern
		parent = parent.Parent()
	}

	return pattern
}

func (r *route) Handler() Handler {
	return r.handler
}

func (r *route) Regexp() *regexp.Regexp {
	return r.regexp
}

func (r *route) Build() {
	matcher := variable.ReplaceAllString(r.Pattern(), "([^(/|$)]+)")
	r.regexp = regexp.MustCompile("^" + matcher + "$")
}

func createRoute(method, pattern string, fn Handler) Route {
	return &route{method: method, pattern: pattern, handler: fn}
}
