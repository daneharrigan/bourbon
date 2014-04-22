package bourbon

import (
	"github.com/bmizerany/assert"
	"github.com/codegangsta/inject"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddlewareContentTypeHandlerBlank(t *testing.T) {
	rw := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com", nil)

	code, encodeable := ContentTypeHandler(rw, r)
	contentType := "application/json; charset=utf-8"

	assert.Equal(t, 0, code)
	assert.Equal(t, nil, encodeable)
	assert.Equal(t, contentType, rw.Header().Get("Content-Type"))
}

func TestMiddlewareContentTypeHandlerWithJSON(t *testing.T) {
	contentType := "application/json; charset=utf-8"
	rw := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	r.Header.Set("Content-Type", contentType)

	code, encodeable := ContentTypeHandler(rw, r)

	assert.Equal(t, 0, code)
	assert.Equal(t, nil, encodeable)
	assert.Equal(t, contentType, rw.Header().Get("Content-Type"))
}

func TestMiddlewareContentTypeHandlerWithVndJSON(t *testing.T) {
	rw := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	r.Header.Set("Content-Type", "application/vnd.example+json; param=2")

	code, encodeable := ContentTypeHandler(rw, r)
	contentType := "application/json; charset=utf-8"

	assert.Equal(t, 0, code)
	assert.Equal(t, nil, encodeable)
	assert.Equal(t, contentType, rw.Header().Get("Content-Type"))
}

func TestMiddlewareContentTypeHandlerWithHTML(t *testing.T) {
	rw := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	r.Header.Set("Content-Type", "text/html")

	code, encodeable := ContentTypeHandler(rw, r)
	contentType := "application/json; charset=utf-8"
	message := Message{
		Code:    415,
		Message: http.StatusText(415),
		Errors:  []string{`"text/html" is not a supported Content-Type`},
	}

	assert.Equal(t, 415, code)
	assert.Equal(t, message, encodeable)
	assert.Equal(t, contentType, rw.Header().Get("Content-Type"))
}

func TestMiddlewareDecodeHandlerWithJSON(t *testing.T) {
	type example struct {
		Name string `json:"name"`
	}
	reader := strings.NewReader(`{"name":"example"}`)
	fn := func(e example) {}

	r, _ := http.NewRequest("GET", "http://example.com", reader)
	c := createTestContext(fn, r)

	code, encodeable := DecodeHandler(c, r)

	assert.Equal(t, 0, code)
	assert.Equal(t, nil, encodeable)
}

func TestMiddlewareDecodeHandlerWithMalformedJSON(t *testing.T) {
	type example struct {
		Name string `json:"name"`
	}
	reader := strings.NewReader(`{"name":""boom}`) // incorrect quoting
	fn := func(e example) {}

	r, _ := http.NewRequest("GET", "http://example.com", reader)
	c := createTestContext(fn, r)

	message := Message{
		Code:    400,
		Message: http.StatusText(400),
		Errors:  []string{"invalid character 'b' after object key:value pair"},
	}

	code, encodeable := DecodeHandler(c, r)

	assert.Equal(t, 400, code)
	assert.Equal(t, message, encodeable)
}

func TestMiddlewareDecodeHandlerWithBlankBody(t *testing.T) {
	type example struct {
		Name string `json:"name"`
	}
	reader := strings.NewReader("") // empty string
	fn := func(e example) {}

	r, _ := http.NewRequest("GET", "http://example.com", reader)
	c := createTestContext(fn, r)

	code, encodeable := DecodeHandler(c, r)

	assert.Equal(t, 0, code)
	assert.Equal(t, nil, encodeable)
}

func TestMiddlewareDecodeHandlerWithContentLengthZero(t *testing.T) {
	fn := func() {}
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	r.Header.Add("Content-Length", "0")
	c := createTestContext(fn, r)

	code, encodeable := DecodeHandler(c, r)

	assert.Equal(t, 0, code)
	assert.Equal(t, nil, encodeable)
}

func createTestContext(fn Handler, r *http.Request) context {
	rt := &Route{new(bourbon), "GET", "/", fn, nil}
	return context{inject.New(), fn, nil, rt, nil, r}
}
