package stdout

import (
	"os"
	"sync"
)

type Stdout struct {
	fd       *os.File
	sendChan chan []byte
	stopOnce sync.Once
}

func New() *Stdout {
	return &Stdout{}
}

func (s *Stdout) Init() error {
	s.sendChan = make(chan []byte)
	s.fd = os.Stdout
	return nil
}

func (s *Stdout) Start() {
	s.start()
}

func (s *Stdout) start() {
	for data := range s.sendChan {
		s.send(data)
	}
}

func (s *Stdout) send(data []byte) {
	_, _ = s.fd.Write(data)
	_, _ = s.fd.Write([]byte("\n"))
}

func (s *Stdout) Stop() {
	s.stopOnce.Do(func() {
		close(s.sendChan)
	})
}

func (s *Stdout) Name() string {
	return "stdout"
}

func (s *Stdout) SendChannel() chan<- []byte {
	return s.sendChan
}
