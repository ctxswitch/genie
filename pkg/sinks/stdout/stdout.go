package stdout

import (
	"os"
)

type Stdout struct {
	fd *os.File
}

func FromConfig() (*Stdout, error) {
	return &Stdout{}, nil
}

func (s *Stdout) Connect() {
	s.fd = os.Stdout
}

func (s *Stdout) Init() {}

func (s *Stdout) Send(data []byte) error {
	_, _ = s.fd.Write(data)
	_, _ = s.fd.Write([]byte("\n"))
	return nil
}

func (s *Stdout) Name() string {
	return "stdout"
}
