package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Middleware receives a content, do whatever it wants it it and returns the
// same content or a new one, as it pleases baby
type MiddleWare func([]byte) ([]byte, error)

// Server is our web server
type Server struct {
	endpoint string
	mdws     []MiddleWare
	addr     string
	port     int
}

// New returns a new server setting all options on it
func New(opts ...Option) *Server {
	s := new(Server)
	s.port = 8080
	s.mdws = make([]MiddleWare, 0)
	s.endpoint = "/"
	for _, o := range opts {
		o(s)
	}
	return s
}

// traverseMiddlewares passes messages through all middlewares, one by one. If
// one fails it returns the error. The output of one middleware is the input of
// the next one.
func (s *Server) traverseMiddlewares(content []byte) (err error) {
	for _, m := range s.mdws {
		content, err = m(content)
		if err != nil {
			return err
		}
	}
	return nil
}

// message is called for every new message. We extract the message body and
// pass it to the traverseMiddlewares function
func (s *Server) message(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	log.Printf("new request has arrived")
	if r.Method != "POST" {
		s.writeError(w, errors.New("invalid request method"))
		return
	}

	reqcontent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.writeError(w, err)
		return
	}
	r.Body.Close()

	err = s.traverseMiddlewares(reqcontent)
	if err != nil {
		s.writeError(w, err)
		return
	}

	log.Printf("message processed")
	// we reach to the end, no problems
	w.WriteHeader(http.StatusAccepted)
}

// writeError is an auxiliar function to help writing errors back to the user
func (s *Server) writeError(w http.ResponseWriter, err error) {

	w.WriteHeader(http.StatusInternalServerError)
	e := Error{err.Error()}

	errstr, _ := json.Marshal(e)
	fmt.Fprintf(w, "%s", errstr)
}

// Bind binds and start to accept http messages. Returns only in case of error
func (s *Server) Bind() error {

	bindaddr := fmt.Sprintf("%s:%d", s.addr, s.port)
	log.Printf("welcome to THE webserver\n")
	log.Printf("listening on %s\n", bindaddr)

	http.HandleFunc(s.endpoint, s.message)
	return http.ListenAndServe(
		bindaddr,
		nil,
	)
}
