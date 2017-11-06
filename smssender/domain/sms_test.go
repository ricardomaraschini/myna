package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleSplit(t *testing.T) {

	content := ""
	for i := 0; i < 154; i++ {
		content = fmt.Sprintf("%sa", content)
	}

	s := SMS{
		Message: content,
	}
	m := s.Split()
	assert.Len(t, m, 2)

	s = SMS{}
	m = s.Split()
	assert.Len(t, m, 0)
}
