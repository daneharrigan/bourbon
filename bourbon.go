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
//   b.Run(b)
package bourbon

import "os"

type bourbon struct {
	parent     Bourbon
	prefix     string
	routes     []Route
	middleware []Handler
	children   []Bourbon
}

var (
	port       string
	router     Router
	server     Server
	middleware []Handler
)

func init() {
	port = os.Getenv("PORT")
	router = createDefaultRouter()
	server = new(defaultServer)
	middleware = append(middleware, ContentTypeHandler, DecodeHandler)
}

// New allocates a new Bourbon.
func New() Bourbon {
	return new(bourbon)
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

func (b *bourbon) SetParent(parent Bourbon) {
	b.parent = parent
}

func (b *bourbon) Parent() Bourbon {
	return b.parent
}

func (b *bourbon) Children() []Bourbon {
	return b.children
}

func (b *bourbon) Mount(child Bourbon) {
	child.SetParent(b)
	b.children = append(b.children, child)
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
	if b.Parent() == nil {
		return b.middleware
	}

	return append(b.Parent().Middleware(), b.middleware)
}

func (b *bourbon) Use(handlers ...Handler) {
	b.middleware = append(b.middleware, handlers...)
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

func (b *bourbon) Run() {
	appendRoute(server.Router(), b)
	server.Run()
}

func (b *bourbon) addRoute(method, pattern string, fn Handler) {
	r := createRoute(method, pattern, fn)
	r.SetParent(b)
	b.routes = append(b.routes, r)
}

func appendRoute(r Router, b Bourbon) {
	r.Add(b.Routes()...)

	for _, child := range b.Children() {
		appendRoute(r, child)
	}
}
