package bourbon

import "net/http"

// Message is a struct for generating consistent error messaging.
type Message struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

// CreateMessage is a convenience function for generating Message structs.
func CreateMessage(code int, errs ...string) Message {
	return Message{Code: code, Message: http.StatusText(code), Errors: errs}
}
