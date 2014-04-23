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
	route      Route
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

func (c *context) Run() {

	defer c.r.Body.Close()
	for _, middleware := range c.middleware {
		c.handleReturns(c.Invoke(middleware))
		if c.rw.Written() {
			return
		}
	}

	c.handleReturns(c.Invoke(c.handler))
}

func createContext(r Route, w http.ResponseWriter, req *http.Request) *context {
	middleware := append(middleware, r.Middleware()...)
	rw := createResponseWriter(w)
	c := &context{inject.New(), r.Handler(), middleware, r, rw, req}
	c.MapTo(c.rw, (*ResponseWriter)(nil))
	c.MapTo(c.rw, (*http.ResponseWriter)(nil))
	c.Map(createParams(c))
	c.Map(r)
	c.Map(c)

	return c
}
