package bourbon_test

import "github.com/daneharrigan/bourbon"

func ExampleRun() {
	b1 := bourbon.New()
	b1.Get("/b1", func(){})

	b2 := bourbon.New()
	b2.Get("/b2", func(){})

	bourbon.Run(b1, b2)
}

func ExampleParams() {
	b := bourbon.New()
	b.Get("/resources/{id}", func(params bourbon.Params){
		println(params["id"])
	})

	bourbon.Run(b)
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
	b.Post("/messages", func(m Message){
		println(m.Value)
	})
}
