package smssender

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/ricardomaraschini/myna/smssender/domain"
	"github.com/stretchr/testify/assert"
)

type clientmock struct {
	msgs int
}

func (c *clientmock) NewMessage(orig string, recpt []string, msg string, par *messagebird.MessageParams) (*messagebird.Message, error) {

	ret := new(messagebird.Message)
	if orig == "genfail" {
		ret.Errors = make([]messagebird.Error, 1)
		ret.Errors[0] = messagebird.Error{}
		return ret, errors.New("an api failure")
	}

	c.msgs++
	return ret, nil
}

func TestSender(t *testing.T) {

	cmock := new(clientmock)
	ch := make(chan []byte)
	s := New(ch)
	s.client = cmock

	ch <- []byte(`
		{
			"originator": "myself",
			"recipient":  "+123412341234",
			"message":    "sample message"
		}
	`)

	time.Sleep(time.Second)
	assert.Equal(t, cmock.msgs, 1)

	content := ""
	for i := 0; i < 161; i++ {
		content = fmt.Sprintf("%sa", content)
	}

	encsms := domain.SMS{
		Message:   content,
		Recipient: "+123123123123",
	}
	jstr, err := json.Marshal(encsms)
	assert.Nil(t, err)

	ch <- jstr

	time.Sleep(time.Second)
	assert.Equal(t, cmock.msgs, 3)
}

func TestSendHugeSMS(t *testing.T) {

	cmock := new(clientmock)
	ch := make(chan []byte)
	s := New(ch)
	s.client = cmock

	content := ""
	for i := 0; i < 100000; i++ {
		content = fmt.Sprintf("%sa", content)
	}

	encsms := domain.SMS{
		Message:   content,
		Recipient: "+123123123123",
	}
	jstr, err := json.Marshal(encsms)
	assert.Nil(t, err)

	ch <- jstr

	time.Sleep(time.Second)
	assert.Equal(t, cmock.msgs, 0)
}

func TestSendWithoutRecipient(t *testing.T) {

	cmock := new(clientmock)
	ch := make(chan []byte)
	s := New(ch)
	s.client = cmock

	ch <- []byte(`
		{
			"originator": "myself",
			"message":    ""
		}
	`)

	time.Sleep(time.Second)
	assert.Equal(t, cmock.msgs, 0)
}

func TestSendEmptySMS(t *testing.T) {

	cmock := new(clientmock)
	ch := make(chan []byte)
	s := New(ch)
	s.client = cmock

	ch <- []byte(`
		{
			"originator": "myself",
			"recipient":  "+21312124123",
			"message":    ""
		}
	`)

	time.Sleep(time.Second)
	assert.Equal(t, cmock.msgs, 0)
}

func TestSendBadSMS(t *testing.T) {

	cmock := new(clientmock)
	ch := make(chan []byte)
	s := New(ch)
	s.client = cmock

	ch <- []byte(`<-as0asdf0->`)

	time.Sleep(time.Second)
	assert.Equal(t, cmock.msgs, 0)
}

func TestApiFailure(t *testing.T) {

	cmock := new(clientmock)
	ch := make(chan []byte)
	s := New(ch)
	s.client = cmock

	ch <- []byte(`
		{
			"originator": "genfail",
			"recipient":  "+21312124123",
			"message":    "a message"
		}
	`)

	time.Sleep(time.Second)
	assert.Equal(t, cmock.msgs, 0)
}
