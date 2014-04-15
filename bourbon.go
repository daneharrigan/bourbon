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

type bourbon struct {
	prefix     string
	routes     []Route
	middleware []Handler
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
	config = createConfig()
	s := config.Server

	for _, b := range bourbons {
		s.Router().Add(b.Routes()...)
	}

	s.Run()
}

// SetConfig accepts a Config struct to override the default Server,
// Router and port number.
//
//   bourbon.SetConfig(Config{
//     Router: new(myCustomRouter)
//   })
func SetConfig(c Config) {
	config = &c
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
	r := createConfig().CreateRoute(method, pattern, fn)
	r.SetParent(b)
	b.routes = append(b.routes, r)
}
