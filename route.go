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

func (r *route) SetParent(parent Bourbon) {
	r.parent = parent
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
	if r.regexp == nil {
		r.createRegexp()
	}
	return r.regexp
}

func (r *route) createRegexp() {
	uri := r.pattern
	if r.parent != nil {
		uri = r.parent.Prefix() + uri
	}
	matchStr := variable.ReplaceAllString(uri, "([^(/|$)]+)")
	r.regexp = regexp.MustCompile("^" + matchStr + "$")
}

func createRoute(method, pattern string, fn Handler) Route {
	return &route{method: method, pattern: pattern, handler: fn}
}
