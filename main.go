package main

import (
	"io/ioutil"

	"github.com/ricardomaraschini/myna/middlewares/jsonvalidation"
	"github.com/ricardomaraschini/myna/middlewares/store"
	"github.com/ricardomaraschini/myna/queues/buffer"
	"github.com/ricardomaraschini/myna/smssender"
	"github.com/ricardomaraschini/myna/webserver"
)

func main() {

	// using memory buffer, if the process dies we lost all pending sms
	buf := buffer.New()
	smssender.New(buf.DrainChannel())

	// the commented lines below make the system use activemq instead
	// of a memory buffer. My implementation of the active mq is a very
	// simple one, no authentication, no special lasers and its presence
	// here exists because a real system running on a docker container
	// must(or at very least should) not persist data.
	/*
		buf, err := activemq.New(
			activemq.WithQueueName("myna"),
		)
		if err != nil {
			panic(err)
		}
		bufch, err := buf.DrainChannel()
		if err != nil {
			panic(err)
		}
		smssender.New(bufch)
	*/

	// loads our json schema file, and here i should have used go generate
	jschema, err := ioutil.ReadFile("assets/sms.json")
	if err != nil {
		panic(err)
	}

	// sets endpoint and creates a chain of middlewares. I also implemented
	// an middleware called `tochannel` where we could easily plug our
	// input into some analytics, I am not using it
	webserver.New(
		webserver.WithEndpoint(
			"/message",
		),
		webserver.WithMiddleware(
			jsonvalidation.New(jschema),
		),
		webserver.WithMiddleware(
			store.New(buf),
		),
	).Bind()
}
