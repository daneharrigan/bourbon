package bourbon

import "regexp"

var variable = regexp.MustCompile(`{([^}]+)}`)

type route struct {
	parent  Bourbon
	method  string
	pattern string
	handler Handler
	params  []string
	matcher *regexp.Regexp
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
	return r.pattern
}

func (r *route) Handler() Handler {
	return r.handler
}

func (r *route) Params() []string {
	return r.params
}

func (r *route) MatchString(uri string) bool {
	if r.matcher == nil {
		r.createMatcher()
	}

	return r.matcher.MatchString(uri)
}

func (r *route) createMatcher() {
	uri := r.parent.Prefix() + r.pattern
	matchStr := variable.ReplaceAllString(uri, "([^(/|$)]+)")
	params := variable.FindAllString(uri, -1)
	for i := 0; i < len(params); i++ {
		params[i] = params[i][1 : len(params[i])-1]
	}

	r.matcher = regexp.MustCompile("^" + matchStr + "$")
	r.params = params
}
