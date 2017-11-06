package tochannel

import "github.com/ricardomaraschini/myna/webserver"

// New returns a simple middleware that dumps all its input to a channel
func New(ch chan []byte) webserver.MiddleWare {
	return func(content []byte) ([]byte, error) {
		ch <- content
		return content, nil
	}
}
