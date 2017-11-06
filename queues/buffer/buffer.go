package buffer

import (
	"sync"
)

// Buffer keeps an slice of []byte. We accept new records from a channel and
// dump the buffer content to another channel. The buffer is operated by two go
// routines, one that keeps reading from input and adding the received records
// into the buffer and the other keeps writing records from the buffer to the
// output. If in a given moment the buffer is empty the go routine that outputs
// records locks until a new record arrives into the buffer. This buffer can
// grow ad eternum and for production release this buffer should keep into
// disk so we can resume sending in case of someone kill this program. Please
// see activemq implementation.
type Buffer struct {
	oncef sync.Once
	onced sync.Once
	mtx   sync.Mutex
	buf   [][]byte
	fill  chan []byte
	drain chan []byte
	ddone chan bool
}

// New returns a new Buffer reference
func New() *Buffer {
	b := new(Buffer)
	b.buf = make([][]byte, 0)
	b.fill = make(chan []byte)
	b.drain = make(chan []byte)
	b.ddone = make(chan bool)
	return b
}

// FillChannel returns the channel from where this buffer is waiting for
// records. It also starts the go routine responsible for reading records from
// the channel
func (b *Buffer) FillChannel() chan []byte {
	b.oncef.Do(func() {
		go func() {
			for {
				b.Add(<-b.fill)
			}
		}()
	})
	return b.fill
}

// Add adds a new record into the buffer
func (b *Buffer) Add(rec []byte) error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.buf = append(b.buf, rec)

	// unlocks drain goroutine
	select {
	case <-b.ddone:
	default:
	}

	return nil
}

// Remove removes first element from the buffer and returns it. Returns
// nil if buffer is empty
func (b *Buffer) Remove() []byte {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	// buffer is empty, blocks
	if len(b.buf) == 0 {
		return nil
	}

	var rec []byte
	rec, b.buf = b.buf[0], b.buf[1:]
	return rec
}

// DrainChannel keeps writing records from the buffer to the output channel. If
// the buffer is empty(drain is faster than fill) it blocks until a new record
// arrives
func (b *Buffer) DrainChannel() chan []byte {

	b.onced.Do(func() {
		go func() {
			for {
				rec := b.Remove()
				if rec == nil {
					// blocks if buffer is empty
					b.ddone <- true
					continue
				}

				b.drain <- rec
			}
		}()
	})

	return b.drain
}
