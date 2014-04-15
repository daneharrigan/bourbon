[![Build Status](https://travis-ci.org/daneharrigan/bourbon.svg)](https://travis-ci.org/daneharrigan/bourbon)
[![GoDoc](https://godoc.org/github.com/daneharrigan/bourbon?status.png)](https://godoc.org/github.com/daneharrigan/bourbon)

# bourbon

Bourbon is a package for rapidly developing JSON web services.

### Getting Started

Start by installing Bourbon:

```console
$ go get github.com/daneharrigan/bourbon
```

Next, we'll create our first Bourbon web service. Let's create a file call
`web.go` containing the following:

```go
package main

import "github.com/daneharrigan/bourbon"

type Example struct {
	Message string
}

func main() {
	b := bourbon.New()
	b.Get("/", func() bourbon.Encodeable {
		return Example{Message: "Hello World!"}
	})

	bourbon.Run(b)
}
```

Lets run the web service locally:

```console
go run web.go
```

Lastly, lets make a request to the web service:

```console
$ curl http://localhost:5000

{"Message":"Hello World!"}
```

Bourbon defaults to port 5000. This can be overwritten by setting the
environment variable `PORT` to the desired value or specifying the value in
`SetPort`.

```go
bourbon.SetPort("3000")
```

### Features

Bourbon's feature set targets JSON web services and aims to make them as easy to
write as possible. Common or repetitive tasks should be handled by Bourbon.

#### Composablity

Bourbon's `Run` function accepts many instances of Bourbon and combines them to
make a single web service. This allows you to extract components into
stand-alone packages and sub-packages.

```go
package main

import (
	"github.com/daneharrigan/bourbon"
	"github.com/example/oauth2"
	"github.com/example/myproject/v1"
	"github.com/example/myproject/v2"
)

func main() {
	v1.API.SetPrefix("/v1")
	v2.API.SetPrefix("/v2")
	bourbon.Run(oauth2.Bourbon, v1.API, v2.API)
}
```

The `oauth2`, `v1` and `v2` packages in this example could be run independently
of each other, but they can also be combined and run as a single web service.

#### Encoding Responses

Bourbon will automatically encode a data structure and write the response when a
`bourbon.Encodeable` is returned by a `Handler`.

```go
package main

import "github.com/daneharrigan/bourbon"

type Example struct {
	Message string
}

func main() {
	b := bourbon.New()
	b.Get("/", func() bourbon.Encodeable {
		return Example{Message: "Hello World!"}
	})

	bourbon.Run(b)
}
```

#### Decoding Requests

When Bourbon sees a data type in the `Handler`'s argument list that does not
belong to the packages `net/http` or `bourbon`, it will decode the request body
into a value of that argument type.

```go
package main

import "github.com/daneharrigan/bourbon"

type Example struct {
	Message string
}

func main() {
	b := bourbon.New()
	b.Post("/", func(e Example) (int, bourbon.Encodeable) {
		println(e.Message)
		return 201, e
	})

	bourbon.Run(b)
}
```

#### HTTP OPTIONS Method

Since you have already declared the URLs and HTTP methods for your web service,
Bourbon will use that information and handle `OPTIONS` requests automatically.

```go
package main

import "github.com/daneharrigan/bourbon"

func main() {
	b := bourbon.New()
	b.Get("/", func(){ /* some work */ })
	b.Put("/", func(){ /* some work */ })
	b.Post("/", func(){ /* some work */ })

	bourbon.Run(b)
}
```

```console
$ curl http://localhost:5000 -X OPTIONS -v

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Content-Length: 0
Allow: GET, PUT, POST
```

#### Reject Non-JSON Requests

Bourbon only accepts JSON. Because of this, Bourbon will reject requests whose
Content-Types are not a form of JSON.

```go
package main

import "github.com/daneharrigan/bourbon"

type Example struct {
	Message string
}

func main() {
	b := bourbon.New()
	b.Get("/", func() bourbon.Encodeable {
		return Example{Message: "Hello World!"}
	})

	bourbon.Run(b)
}
```

With the Content-Type of `application/json`:

```console
$ curl http://localhost:5000 -H "Content-Type: application/json"
{"Message":"Hello World!"}
```

With the Content-Type of `application/vnd.custom+json`:

```console
$ curl http://localhost:5000 -H "Content-Type: application/vnd.custom+json"
{"Message":"Hello World!"}
```

With the Content-Type of `text/html`:

```console
$ curl http://localhost:5000 -H "Content-Type: text/html"
{"code":415,"message":"Unsupported Media Type","errors":["\"text/html\" is not a supported Content-Type"]}
```

### Need Help

Did I overlook something? Is an area of Bourbon confusing? Or maybe you just
want to say hi! Feel free to reach out via Twitter, email or over Github.

- [@daneharrigan](https://twitter.com/daneharrigan)
- [dane@heroku.com](mailto:dane@heroku.com)
- [Github Issue](https://github.com/daneharrigan/bourbon/issues)
