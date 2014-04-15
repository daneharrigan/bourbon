package bourbon

import "net/http"

type server struct {
}

func (s *server) Run() {
	if defaultPort == "" {
		defaultPort = "5000"
	}
	http.ListenAndServe(":"+defaultPort, s)
}

func (s *server) Router() Router {
	return defaultRouter
}

func (s *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	action := s.Router().Find(r.Method, r.URL.Path)
	action.Run(rw, r)
}
