package store

import "github.com/ricardomaraschini/myna/webserver"

// Repository is an entity capable of saving an []byte somewhere. As a side
// note both memory buffer and active mq implement this interface.
type Repository interface {
	Add([]byte) error
}

// New returns a simple middleware that save content into a repository
func New(rep Repository) webserver.MiddleWare {
	return func(content []byte) ([]byte, error) {
		return content, rep.Add(content)
	}
}
