package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Rep struct {
	records [][]byte
}

func (r *Rep) Add(rec []byte) error {
	r.records = append(r.records, rec)
	return nil
}

func TestStore(t *testing.T) {
	r := new(Rep)
	middleware := New(r)
	for i := 0; i < 100; i++ {
		middleware([]byte("123"))
	}

	assert.Len(t, r.records, 100)
}
