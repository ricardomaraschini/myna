package webserver

type Option func(*Server)

// WithAddress sets the ip address where our server is going to bind to
func WithAddress(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithPort sets the port where our server is going to bind to
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// WithMiddleware adds a middleware to the server's list of middlewares
func WithMiddleware(m MiddleWare) Option {
	return func(s *Server) {
		s.mdws = append(s.mdws, m)
	}
}

// WithEndpoint sets the endpoint
func WithEndpoint(ep string) Option {
	return func(s *Server) {
		s.endpoint = ep
	}
}
