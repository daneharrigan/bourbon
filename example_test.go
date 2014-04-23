package bourbon_test

import (
	"github.com/daneharrigan/bourbon"
	"net/http"
)

func ExampleRun() {
	b := bourbon.New()
	b.Get("/", func() {})
	b.Run()
}

func ExampleParams() {
	b := bourbon.New()
	b.Get("/resources/{id}", func(params bourbon.Params) {
		println(params["id"])
	})

	b.Run()
}

func ExampleEncodeable() {
	b := bourbon.New()
	b.Get("/", func() bourbon.Encodeable {
		var e struct {
			Message string
		}

		e.Message = "Hello World!"
		return e
	})

}

func ExampleDecodeHandler() {
	type Message struct {
		Value string
	}

	b := bourbon.New()
	b.Post("/messages", func(m Message) {
		println(m.Value)
	})
}

func ExampleHandler() {
	b := bourbon.New()
	// use existing net/http handlers
	b.Get("/legacy", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(`{"message": "hello world"}`))
	})

	// return an int as the status code
	b.Get("/example/1", func(rw http.ResponseWriter) int {
		rw.Write([]byte(`{"message": "hello world"}`))
		return 200
	})

	// return an int as the status code and an encodeable data structure
	b.Get("/example/2", func() (int, bourbon.Encodeable) {
		var item struct {
			Value int
		}

		return 200, item
	})

	// decode request body
	// POST /example
	// { "Name": "Test" }

	type Example struct {
		Name string
	}

	// return an int as the status code and an encodeable data structure
	// accept a custom data structure, Example in this case, and Bourbon
	// will decode the request body into it.
	b.Post("/example", func(e Example) (int, bourbon.Encodeable) {
		return 201, e
	})
}

func ExampleMount() {
	// The child Bourbon will inherit middleware and URL prefixes from all
	// of it's ascending Bourbon parents.
	parent := bourbon.New()
	first := bourbon.New()
	last := bourbon.New()

	parent.Mount(first)
	first.Mount(last)
	parent.Run()
}
