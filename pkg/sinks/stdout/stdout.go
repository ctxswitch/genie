package stdout

import (
	"context"
	"os"
	"sync"
)

type Stdout struct {
	fd       *os.File
	sendChan chan []byte
	stopChan chan struct{}
	stopOnce sync.Once
}

func New() *Stdout {
	return &Stdout{}
}

func (s *Stdout) Init() error {
	s.sendChan = make(chan []byte, 100)
	s.stopChan = make(chan struct{})
	s.fd = os.Stdout
	return nil
}

func (s *Stdout) Start(ctx context.Context) {
	s.start(ctx)
}

func (s *Stdout) start(ctx context.Context) {
	for {
		select {
		case <-s.stopChan:
			return
		case <-ctx.Done():
			return
		case d := <-s.sendChan:
			s.send(d)
		}
	}
}

func (s *Stdout) send(data []byte) error {
	_, _ = s.fd.Write(data)
	_, _ = s.fd.Write([]byte("\n"))
	return nil
}

func (s *Stdout) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopChan)
	})
}

func (s *Stdout) Name() string {
	return "stdout"
}

func (s *Stdout) SendChannel() chan<- []byte {
	return s.sendChan
}
