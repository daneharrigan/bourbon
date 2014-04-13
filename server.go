package bourbon

import "net/http"

type server struct {
}

func (s *server) Run() {
	http.ListenAndServe(":"+config.Port, s)
}

func (s *server) Router() Router {
	return config.Router
}

func (s *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	action := s.Router().Find(r.Method, r.URL.Path)
	action.Run(rw, r)
}
