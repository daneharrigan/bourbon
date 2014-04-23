package bourbon

import (
	"net/http"
	"regexp"
)

// Encodeable is any data structure that encodes to JSON.
type Encodeable interface{}

// Handler is any callable function. A Handler can accept a range of arguments
// from the packages net/http and bourbon which will be injected automatically.
//
// When a struct is found in the argument list that does not belong to the
// packages bourbon or net/http, the request body will automatically be decoded
// into the struct and passed into the function.
//
// Handlers can return zero, one, or two values. If an integer is returned, it is
// usd as the status code. If an Encodeable is returned, it is encoded and
// written to the response.
type Handler interface{}

// Server is Bourbon's Server interface. It accepts and coordinates requests
// with the router to find the appropriate response.
type Server interface {
	// Router returns the Router containing every route.
	Router() Router

	// Run executes the Server and binds to the port declared in Config's
	// Port.
	Run()

	// ServeHTTP is the entry point from net/http into a bourbon Server.
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// Router is Bourbon's Router interface. It tracks every route defined in a
// Bourbon.
type Router interface {
	// Add appends routes to the Router.
	Add(...Route)

	// Find accepts the request method, URL and returns an Action.
	Find(string, string) Action
}

// Route is Bourbon's Route interface. It stores the HTTP request method, URL
// pattern, Handler and Bourbon parent on which it was declared. A Route is
// appened to the Router's list of routes and queried on each HTTP request.
type Route interface {
	// SetParent declares which Bourbon owns the Route.
	SetParent(Bourbon)

	// Parent returns the Bourbon that owns the Route. The parent is
	// referenced if the Route has a URL prefix and any middleware that
	// should be called when handling the request.
	Parent() Bourbon

	// Method returns the HTTP request method.
	Method() string

	// Pattern returns the URL pattern used to match against URLs. The
	// pattern may contain parameters declared with the `{name}` syntax.
	//
	//   r.Pattern() // => /resources/{resource_id}/example
	Pattern() string

	// Handler returns the function with which to process the incoming
	// request.
	Handler() Handler

	// Regexp returns the regular expression used for matching the route
	// against the request URL. It is also used to read parameters from the
	// URL.
	Regexp() *regexp.Regexp

	// Build creates and assigns the regular expression used to match the
	// request URL and read parameters from the URL.
	Build()
}

// Action is Bourbon's interface for responding to a request.
type Action interface {
	// Run invokes the middleware and Handler scoped to the request.
	Run(http.ResponseWriter, *http.Request)
}

// Bourbon is the initial interface in the Bourbon package.
type Bourbon interface {
	// Parent returns the parent Bourbon structure
	Parent() Bourbon

	// SetPrefix accepts a string to prefix every route in the Bourbon.
	//
	//   v1 := bourbon.New()
	//   v1.SetPrefix("/v1")
	SetPrefix(string)

	// Prefix returns the string used prefix every route in the Bourbon.
	Prefix() string

	// Use appends middleware to be used on each request in the Bourbon.
	// Middleware is scoped to a single Bourbon. Running mutliple Bourbons
	// at the same time will not combine middleware.
	//
	//   public  := bourbon.New() // does not use basic auth
	//   private := bourbon.New() // does use basic auth
	//   private.Use(BasicAuthHandler)
	//   bourbon.Run(public, private)
	Use(...Handler)

	// Middleware returns all of the middleware used by the Bourbon.
	Middleware() []Handler

	// Get declares route within the Bourbon that responds to HTTP's GET
	// method and the URL pattern provided.
	//
	//   b := bourbon.New()
	//   b.Get("/resources/{id}", func() {})
	Get(string, Handler)

	// Put declares route within the Bourbon that responds to HTTP's PUT
	// method and the URL pattern provided.
	//
	//   b := bourbon.New()
	//   b.Put("/resources/{id}", func() {})
	Put(string, Handler)

	// Post declares route within the Bourbon that responds to HTTP's POST
	// method and the URL pattern provided.
	//
	//   b := bourbon.New()
	//   b.Put("/resources/{id}", func() {})
	Post(string, Handler)

	// Head declares route within the Bourbon that responds to HTTP's HEAD
	// method and the URL pattern provided.
	//
	//   b := bourbon.New()
	//   b.Post("/resources/{id}", func() {})
	Head(string, Handler)

	// Patch declares route within the Bourbon that responds to HTTP's PATCH
	// method and the URL pattern provided.
	//
	//   b := bourbon.New()
	//   b.Patch("/resources/{id}", func() {})
	Patch(string, Handler)

	// Delete declares route within the Bourbon that responds to HTTP's
	// DELETE method and the URL pattern provided.
	//
	//   b := bourbon.New()
	//   b.Delete("/resources/{id}", func() {})
	Delete(string, Handler)

	// Routes returns a slice of routes defined on the Bourbon
	Routes() []Route
}

// ResponseWriter is Bourbon's interface for responding to HTTP requests.
type ResponseWriter interface {
	// Write writes a slice of bytes to the HTTP response.
	Write([]byte) (int, error)

	// WriteHeader sets the status code of the HTTP response.
	WriteHeader(int)

	// Header returns http.Header.
	Header() http.Header

	// Stream accepts encodeable data types and writes them to the response.
	// This function can be called multiple times to stream JSON  chunks to
	// the response.
	Stream(Encodeable)

	// Written returns a boolean indicating whether or not data has been
	// written to the response.
	Written() bool
}
