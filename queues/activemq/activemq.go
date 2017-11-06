package activemq

import (
	"fmt"
	"time"

	"github.com/go-stomp/stomp"
)

// ActiveMQ operates as a buffer controller. It push and pull messages from
// an activemq queue
type ActiveMQ struct {
	conn  *stomp.Conn
	addr  string
	port  int
	qname string
	drain chan []byte
}

// New returns a new ActiveMQ reference, default server is set to localhost,
// default port is 61613 and default queue name is "test". We use stomp to
// connect to the remote server and return an error if connection attempt
// fails
func New(opts ...Option) (*ActiveMQ, error) {
	a := new(ActiveMQ)
	a.addr = "127.0.0.1"
	a.port = 61613
	a.qname = "test"
	a.drain = make(chan []byte)

	for _, o := range opts {
		o(a)
	}

	// XXX
	// by default we set the activemq idle connection timeout to one
	// year, this obviously need to be changed and properly handled but
	// for sake of this prototype it should be enough
	conn, err := stomp.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", a.addr, a.port),
		stomp.ConnOpt.HeartBeatError(210240*time.Hour),
	)
	if err != nil {
		return a, err
	}

	a.conn = conn
	return a, nil
}

// DrainChannel returns a channel to where ActiveMQ writes all new messages
func (a *ActiveMQ) DrainChannel() (chan []byte, error) {

	q := fmt.Sprintf("/queue/%s", a.qname)
	sub, err := a.conn.Subscribe(q, stomp.AckAuto)
	if err != nil {
		return a.drain, err
	}

	go func() {
		for {
			msg := <-sub.C
			if msg != nil {
				a.drain <- msg.Body
			}
		}
	}()

	return a.drain, nil
}

// Add enqueues a new message on ActiveMQ queue
func (a *ActiveMQ) Add(rec []byte) error {

	q := fmt.Sprintf("/queue/%s", a.qname)
	return a.conn.Send(
		q,
		"text/plain",
		rec,
		stomp.SendOpt.Receipt,
	)
}
