package bourbon

import "net/http"

type defaultServer struct {
}

func (s *defaultServer) Run() {
	if port == "" {
		port = "5000"
	}
	http.ListenAndServe(":"+port, s)
}

func (s *defaultServer) Router() Router {
	return router
}

func (s *defaultServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r := s.Router().Find(req.Method, req.URL.Path)
	c := createContext(r, rw, req)
	c.Run()
}
