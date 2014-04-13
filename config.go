package bourbon

import "os"

var config *Config

// Config stores the configuration of Bourbon and it's components. The Server,
// Router, Port number and function for creating new Routes are declared in
// Config. Each of these components can be overwritten by creating a custom
// Config and passing it into the Bourbon package.
//
//   // overwriting the Router
//   bourbon.SetConfig(Config{Router: new(myCustomRouter)})
type Config struct {
	Server      Server
	Router      Router
	Port        string
	CreateRoute func(method, pattern string, fn Handler) Route
}

func createConfig() *Config {
	if config == nil {
		config = new(Config)
	}

	if config.Server == nil {
		config.Server = new(server)
	}

	if config.Router == nil {
		config.Router = &router{routes: make(map[string][]Route)}
	}

	if config.Port == "" {
		config.Port = os.Getenv("PORT")
		if config.Port == "" {
			config.Port = "5000"
		}
	}

	if config.CreateRoute == nil {
		config.CreateRoute = defaultCreateRoute
	}

	return config
}

func defaultCreateRoute(method, pattern string, fn Handler) Route {
	return &route{method: method, pattern: pattern, handler: fn}
}
