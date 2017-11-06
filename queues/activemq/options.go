package activemq

// Option sets an option on an ActiveMQ reference
type Option func(*ActiveMQ)

// WithAddress sets active mq server ip address
func WithAddress(addr string) Option {
	return func(a *ActiveMQ) {
		a.addr = addr
	}
}

// WithPort set active mq server tcp port
func WithPort(port int) Option {
	return func(a *ActiveMQ) {
		a.port = port
	}
}

// WithQueueName sets the default queue name
func WithQueueName(qname string) Option {
	return func(a *ActiveMQ) {
		a.qname = qname
	}
}
