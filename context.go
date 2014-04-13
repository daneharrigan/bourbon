package bourbon

import (
	"encoding/json"
	"github.com/codegangsta/inject"
	"net/http"
	"reflect"
)

type context struct {
	inject.Injector
	handler    Handler
	middleware []Handler
	rw         *responseWriter
	r          *http.Request
	encoder    *json.Encoder
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
			c.encoder.Encode(v.Interface())
		default:
			panic("Unexpected return type")
		}
	}
}

func (c *context) Run(rw http.ResponseWriter, r *http.Request) {
	res := &responseWriter{rw, false}
	c.encoder = json.NewEncoder(res)
	c.rw = res
	c.r = r
	c.MapTo(res, (*http.ResponseWriter)(nil))
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

func createContext(r Route) *context {
	c := &context{inject.New(), r.Handler(), nil, nil, nil, nil}
	if r.Parent() != nil {
		c.middleware = r.Parent().Middleware()
	}

	return c
}
