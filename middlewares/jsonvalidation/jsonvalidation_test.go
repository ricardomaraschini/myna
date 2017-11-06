package jsonvalidation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {

	schema := []byte(`{
		"type": "object"
	}`)

	middleware := New(schema)
	content, err := middleware([]byte("123"))
	assert.NotNil(t, err)
	assert.Nil(t, content)

	schema = []byte("")
	middleware = New(schema)
	_, err = middleware([]byte(""))

	schema = []byte(`{
		"type": "object"
	}`)
	middleware = New(schema)
	_, err = middleware([]byte("{}"))
	assert.Nil(t, err)

}
