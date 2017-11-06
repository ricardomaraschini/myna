package tochannel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToChannel(t *testing.T) {

	c := make(chan []byte, 1)
	middleware := New(c)
	_, err := middleware([]byte("testing 123"))
	assert.Nil(t, err)

	res := <-c
	assert.Equal(t, res, []byte("testing 123"))

}
