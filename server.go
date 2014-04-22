package bourbon

import "net/http"

type defaultServer struct{}

func (s *defaultServer) Run() {
	if port == "" {
		port = "5000"
	}
	http.ListenAndServe(":"+port, s)
}

func (s *defaultServer) Router() Router {
	return router
}

func (s *defaultServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	action := s.Router().Find(r.Method, r.URL.Path)
	action.Run(rw, r)
}
