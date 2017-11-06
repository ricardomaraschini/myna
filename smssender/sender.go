package smssender

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/ricardomaraschini/myna/smssender/domain"
)

// MBClient exists only for test mocking
type MBClient interface {
	NewMessage(
		string,
		[]string,
		string,
		*messagebird.MessageParams,
	) (*messagebird.Message, error)
}

type Sender struct {
	in     chan []byte
	client MBClient
	key    string
}

// Sender keeps the needed information to send sms to remote server
func New(readfrom chan []byte) *Sender {
	s := new(Sender)
	s.in = readfrom
	s.client = messagebird.New(os.Getenv("MBKEY"))
	go s.start()
	return s
}

// start initiate the process of draining the buffer, one sms per second
func (s *Sender) start() {
	for ctnt := range s.in {
		msg := domain.SMS{}
		err := json.Unmarshal(ctnt, &msg)
		if err != nil {
			log.Print(err.Error())
			continue
		}

		go s.sendsms(msg)
		time.Sleep(time.Second)
	}
}

// senssms sends a sms to remote endpoint
func (s *Sender) sendsms(sms domain.SMS) {

	// Avoid to call remote api in case of incomplete sms
	if sms.Recipient == "" {
		log.Println("empty recipient")
		return
	}

	parts := sms.Split()
	if len(parts) == 0 {
		log.Println("empty sms message")
		return
	}

	// XXX
	// As we have only one byte for the total os chunks(0xFF) lets
	// avoid to send messages bigger than that so we don't break our
	// naive UDH implementation
	if len(parts) > 255 {
		log.Println("message too big")
		return
	}

	msgid := rand.Intn(255)
	for i := 0; i < len(parts); i++ {

		log.Printf("sending part %d of %d\n", i+1, len(parts))
		msg, err := s.client.NewMessage(
			sms.Originator,
			[]string{sms.Recipient},
			parts[i],
			s.newMsgParams(msgid, len(parts), i+1),
		)
		if err != nil {
			for _, e := range msg.Errors {
				log.Print(e.Description)
			}
			return
		}
	}

	log.Printf("sms sent\n")
}

// newMsgParams creates a new MessageParams setting its internal udh value.
// by now the udh implemetation is very naive, we just set very properties
// for basic multipart messages.
func (s *Sender) newMsgParams(id, total, nr int) *messagebird.MessageParams {
	return &messagebird.MessageParams{
		Type: "binary",
		TypeDetails: messagebird.TypeDetails{
			"udh": fmt.Sprintf(
				"050003%02X%02X%02X", id, total, nr,
			),
		},
	}
}
