package bourbon

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateParams(t *testing.T) {
	r := &Route{Parent: new(bourbon), Pattern: "/one/{one_id}/two/{two_id}"}
	r.Build()
	req, _ := http.NewRequest("GET", "http://example.com/one/1/two/2", nil)
	c := &context{route: r, r: req}
	params := createParams(c)

	assert.Equal(t, "1", params["one_id"])
	assert.Equal(t, "2", params["two_id"])
}

func TestParamsArgument(t *testing.T) {
	fn := func(params Params, rw http.ResponseWriter) {
		rw.Write([]byte(params["one_id"]))
	}

	req, _ := http.NewRequest("GET", "http://example.com/one/1",
		strings.NewReader(""))
	rw := httptest.NewRecorder()
	c := createTestContext(fn, req)
	c.route = &Route{Parent: new(bourbon), Pattern: "/one/{one_id}", Handler: fn}
	c.route.Build()
	c.Run(rw, req)

	body, _ := ioutil.ReadAll(rw.Body)
	assert.Equal(t, []byte("1"), body)
}
