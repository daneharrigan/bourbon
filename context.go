package bourbon

import (
	"github.com/codegangsta/inject"
	"net/http"
	"reflect"
)

type context struct {
	inject.Injector
	handler    Handler
	middleware []Handler
	route      *Route
	rw         ResponseWriter
	r          *http.Request
}

func (c *context) handleReturns(values []reflect.Value, err error) {
	for _, v := range values {
		switch {
		case v.Kind() == reflect.Int:
			if v.Int() > 0 {
				c.rw.WriteHeader(int(v.Int()))
			}
		case v.IsNil():
			continue
		case v.Kind() == reflect.Interface:
			c.rw.Stream(v.Interface())
		default:
			panic("Unexpected return type")
		}
	}
}

func (c *context) Run(rw http.ResponseWriter, r *http.Request) {
	c.rw = createResponseWriter(rw)
	c.r = r
	c.MapTo(c.rw, (*ResponseWriter)(nil))
	c.MapTo(c.rw, (*http.ResponseWriter)(nil))
	c.Map(createParams(c))
	c.Map(r)
	c.Map(c)

	defer c.r.Body.Close()
	for _, middleware := range c.middleware {
		c.handleReturns(c.Invoke(middleware))
		if c.rw.Written() {
			return
		}
	}

	c.handleReturns(c.Invoke(c.handler))
}

func createContext(r *Route) *context {
	c := &context{inject.New(), r.Handler, nil, r, nil, nil}
	if r.Parent != nil {
		c.middleware = r.Parent.Middleware()
	}

	return c
}
