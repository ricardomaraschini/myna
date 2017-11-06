package buffer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {

	b := New()
	out := b.DrainChannel()
	in := b.FillChannel()

	// dispatch a reader from the buffer
	records := 0
	go func() {
		for {
			<-out
			records++
		}
	}()

	// writes 100 records into the buffer
	for i := 0; i < 100; i++ {
		in <- []byte("123")
	}

	// wait a little so the last message is processed
	time.Sleep(time.Millisecond * 10)
	assert.Len(t, b.buf, 0)
	assert.Equal(t, records, 100)

	for i := 0; i < 1000; i++ {
		in <- []byte("321")
	}
	time.Sleep(time.Millisecond * 10)
	assert.Len(t, b.buf, 0)
	assert.Equal(t, records, 1100)
}
