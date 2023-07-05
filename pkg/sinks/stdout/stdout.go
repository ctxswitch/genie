package stdout

import (
	"os"

	"ctx.sh/genie/pkg/sinks"
)

type Stdout struct {
	fd *os.File
}

func (s *Stdout) Connect() {
	s.fd = os.Stdout
}

func (s *Stdout) Init() {}

func (s *Stdout) Send(event []byte) {
	_, _ = s.fd.Write(event)
}

var _ sinks.Sink = &Stdout{}
