package bourbon

import (
	"net/http"
	"strings"
)

type router struct {
	routes map[string][]Route
}

func (rt *router) Add(routes ...Route) {
	for _, r := range routes {
		rt.routes[r.Method()] = append(rt.routes[r.Method()], r)
	}
}

func (rt *router) Find(method, uri string) Action {
	// serve OPTIONS
	if method == "OPTIONS" {
		var methods []string
		var parent Bourbon
		for k, _ := range rt.routes {
			for _, r := range rt.routes[k] {
				if r.MatchString(uri) {
					methods = append(methods, k)
					parent = r.Parent()
				}
			}
		}

		if len(methods) > 0 {
			options := createOptions(methods)
			options.SetParent(parent)
			return createContext(options)
		}

		return createContext(createNotFound())
	}

	// serve route
	for _, r := range rt.routes[method] {
		if r.MatchString(uri) {
			return createContext(r)
		}
	}

	// serve 405
	for m, routes := range rt.routes {
		if m == method {
			continue
		}

		for _, r := range routes {
			if r.MatchString(uri) {
				methodNotAllowed := createMethodNotAllowed()
				methodNotAllowed.SetParent(r.Parent())
				return createContext(methodNotAllowed)
			}
		}
	}

	// serve 404
	return createContext(createNotFound())
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
			return 405, createMessage(405)
		},
	}
}

func createNotFound() Route {
	return &route{
		handler: func() (int, Encodeable) {
			return 404, createMessage(404)
		},
	}
}
