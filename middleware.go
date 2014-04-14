package bourbon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// ContentTypeHandler is a middleware for handling incoming and outgoing
// Content-Types. Incoming Content-Types that are not JSON are rejected.
// Outgoing resposnes are given an application/json Content-Type. This
// middleware automatically prepends routes added to Bourbon.
func ContentTypeHandler(rw http.ResponseWriter, r *http.Request) (int, Encodeable) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	contentType := strings.Split(r.Header.Get("Content-Type"), ";")[0]
	size := len(contentType)

	if size == 0 || (size > 4 && contentType[size-4:] == "json") {
		return 0, nil
	}

	err := fmt.Sprintf("%q is not a supported Content-Type", contentType)
	message := createMessage(415)
	message.Errors = append(message.Errors, err)
	return 415, message
}

// DecodeHandler is a middleware for decoding JSON request bodies into structs.
// The middleware will analyze the argument list of the route's Handler to
// determine if the request body should be decoded. If the argument list
// contains a struct type that does not belong to the net/http or bourbon
// package, DecodeHandler assumes the request body should be decoded into a
// value of that type and passed into the route's Handler.
func DecodeHandler(c context, r *http.Request) (int, Encodeable) {
	if r.ContentLength == 0 {
		return 0, nil
	}

	typeOf := reflect.TypeOf(c.handler)
	for i := 0; i < typeOf.NumIn(); i++ {
		argument := typeOf.In(i)
		if value := c.Get(argument); value.IsValid() {
			continue
		}

		value := reflect.New(argument)
		err := json.NewDecoder(r.Body).Decode(value.Interface())
		if err != nil {
			message := createMessage(400)
			message.Errors = append(message.Errors, err.Error())
			return 400, message
		}

		c.Map(reflect.Indirect(value).Interface())
		break
	}

	return 0, nil
}
