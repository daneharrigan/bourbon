package bourbon

import "net/http"

type Message struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func createMessage(code int) Message {
	return Message{Code: code, Message: http.StatusText(code)}
}
