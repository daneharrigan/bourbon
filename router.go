package bourbon

import (
	"net/http"
	"strings"
)

type defaultRouter struct {
	routes map[string][]Route
}

func (dr *defaultRouter) Add(routes ...Route) {
	for _, r := range routes {
		dr.routes[r.Method()] = append(dr.routes[r.Method()], r)
	}
}

func (dr *defaultRouter) Find(method, uri string) Route {
	// serve OPTIONS
	if method == "OPTIONS" {
		var methods []string
		var parent Bourbon
		for k := range dr.routes {
			for _, r := range dr.routes[k] {
				if r.Regexp().MatchString(uri) {
					methods = append(methods, k)
					parent = r.Parent()
				}
			}
		}

		if len(methods) > 0 {
			options := createOptions(methods)
			options.SetParent(parent)
			return options
		}

		return createNotFound()
	}

	// serve route
	for _, r := range dr.routes[method] {
		if r.Regexp().MatchString(uri) {
			return r
		}
	}

	// serve 405
	for m, routes := range dr.routes {
		if m == method {
			continue
		}

		for _, r := range routes {
			if r.Regexp().MatchString(uri) {
				methodNotAllowed := createMethodNotAllowed()
				methodNotAllowed.SetParent(r.Parent())
				return methodNotAllowed
			}
		}
	}

	// serve 404
	return createNotFound()
}

func createOptions(methods []string) Route {
	return &route{
		handler: func(rw http.ResponseWriter) {
			rw.Header().Set("Allow", strings.Join(methods, ","))
			rw.Header().Set("Content-Length", "0")
		},
	}
}

func createMethodNotAllowed() Route {
	return &route{
		handler: func() (int, Encodeable) {
			return 405, CreateMessage(405)
		},
	}
}

func createNotFound() Route {
	return &route{
		handler: func() (int, Encodeable) {
			return 404, CreateMessage(404)
		},
	}
}

func createDefaultRouter() Router {
	return &defaultRouter{routes: make(map[string][]Route)}
}
