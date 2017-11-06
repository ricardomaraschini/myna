package webserver

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	srv := New()
	assert.Equal(t, srv.addr, "")
	assert.Equal(t, srv.port, 8080)

	srv = New(
		WithPort(1234),
		WithAddress("localhost"),
		WithEndpoint("/test"),
	)
	assert.Equal(t, srv.addr, "localhost")
	assert.Equal(t, srv.port, 1234)
	assert.Equal(t, srv.endpoint, "/test")
}

func TestWithSomeMiddlwares(t *testing.T) {

	valid := func(content []byte) ([]byte, error) {
		assert.Equal(t, content, []byte("123321"))
		return []byte("abc"), nil
	}

	send := func(content []byte) ([]byte, error) {
		assert.Equal(t, content, []byte("abc"))
		return content, nil
	}

	invalid := func(content []byte) ([]byte, error) {
		return content, errors.New("error")
	}

	srv := New(
		WithMiddleware(valid),
		WithMiddleware(send),
	)

	assert.Len(t, srv.mdws, 2)
	err := srv.traverseMiddlewares([]byte("123321"))
	assert.Nil(t, err)

	srv = New(
		WithMiddleware(valid),
		WithMiddleware(send),
		WithMiddleware(invalid),
	)

	err = srv.traverseMiddlewares([]byte("123321"))
	assert.NotNil(t, err)

}
