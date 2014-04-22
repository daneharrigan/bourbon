// Package bourbon is a package for rapidly developing JSON web services.
//
//   b := bourbon.New()
//   b.Get("/", func() bourbon.Encodeable {
//     var resource struct {
//       Name string
//     }
//
//     return resource
//   })
//
//   bourbon.Run(b)
package bourbon

import "os"

type bourbon struct {
	prefix     string
	routes     []Route
	middleware []Handler
}

var (
	port   string
	router Router
	server Server
)

func init() {
	port = os.Getenv("PORT")
	server = new(defaultServer)
	router = &defaultRouter{routes: make(map[string][]Route)}
}

// New allocates a new Bourbon.
func New() Bourbon {
	b := new(bourbon)
	b.Use(ContentTypeHandler)
	b.Use(DecodeHandler)

	return b
}

// Run combines all Bourbons into a Server and runs the server. Use Run with one
// or more Bourbons to keep the API modular and composable.
func Run(bourbons ...Bourbon) {
	for _, b := range bourbons {
		server.Router().Add(b.Routes()...)
	}

	server.Run()
}

// SetRouter accepts a struct that implements that Router interface and replaces
// Bourbon's default router.
func SetRouter(r Router) {
	router = r
}

// SetServer accepts a struct that implements that Server interface and replaces
// Bourbon's default Server.
func SetServer(s Server) {
	server = s
}

// SetPort accepts a port as a string and overrides the default port "5000". The
// default port can also be overwritten by setting the environment variable PORT
// to the desired value.
func SetPort(p string) {
	port = p
}

func (b *bourbon) SetPrefix(prefix string) {
	b.prefix = prefix
}

func (b *bourbon) Prefix() string {
	return b.prefix
}

func (b *bourbon) Routes() []Route {
	return b.routes
}

func (b *bourbon) Middleware() []Handler {
	return b.middleware
}

func (b *bourbon) Use(middleware ...Handler) {
	b.middleware = append(b.middleware, middleware...)
}

func (b *bourbon) Get(pattern string, fn Handler) {
	b.addRoute("GET", pattern, fn)
}

func (b *bourbon) Put(pattern string, fn Handler) {
	b.addRoute("PUT", pattern, fn)
}

func (b *bourbon) Post(pattern string, fn Handler) {
	b.addRoute("POST", pattern, fn)
}

func (b *bourbon) Head(pattern string, fn Handler) {
	b.addRoute("HEAD", pattern, fn)
}

func (b *bourbon) Patch(pattern string, fn Handler) {
	b.addRoute("PATCH", pattern, fn)
}

func (b *bourbon) Delete(pattern string, fn Handler) {
	b.addRoute("DELETE", pattern, fn)
}

func (b *bourbon) addRoute(method, pattern string, fn Handler) {
	r := createRoute(method, pattern, fn)
	r.SetParent(b)
	b.routes = append(b.routes, r)
}
