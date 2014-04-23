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

func (s *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r := s.Router().Find(req.Method, req.URL.Path)
	c := createContext(r, rw, req)
	c.Run()
}
