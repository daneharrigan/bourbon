package bourbon

import "regexp"

var variable = regexp.MustCompile(`{([^}]+)}`)

type Route struct {
	Parent  Bourbon
	Method  string
	Pattern string
	Handler Handler
	Regexp  *regexp.Regexp
}

func (r *Route) Build() {
	uri := r.Pattern
	if r.Parent != nil {
		uri = r.Parent.Prefix() + uri
	}
	matcher := variable.ReplaceAllString(uri, "([^(/|$)]+)")
	r.Regexp = regexp.MustCompile("^" + matcher + "$")
}

func createRoute(method, pattern string, fn Handler) *Route {
	return &Route{Method: method, Pattern: pattern, Handler: fn}
}
