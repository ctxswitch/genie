package stdout

import (
	"os"
	"sync"
)

// Stdout is a sink that writes to stdout.
type Stdout struct {
	fd       *os.File
	sendChan chan []byte
	stopOnce sync.Once
}

// New returns a new stdout sink.
func New() *Stdout {
	return &Stdout{
		sendChan: make(chan []byte),
	}
}

// Init initializes the stdout sink.
func (s *Stdout) Init() error {
	s.fd = os.Stdout
	return nil
}

// Start starts the stdout sink.
func (s *Stdout) Start() {
	s.start()
}

// start starts the stdout sink.
func (s *Stdout) start() {
	for data := range s.sendChan {
		s.send(data)
	}
}

// send writes the data to stdout.
func (s *Stdout) send(data []byte) {
	_, _ = s.fd.Write(data)
	_, _ = s.fd.Write([]byte("\n"))
}

// Stop stops the stdout sink.
func (s *Stdout) Stop() {
	s.stopOnce.Do(func() {
		close(s.sendChan)
	})
}

// SendChannel returns the send channel for the stdout sink.
func (s *Stdout) SendChannel() chan<- []byte {
	return s.sendChan
}
