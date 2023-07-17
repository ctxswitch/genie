package stdout

import (
	"os"

	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/sinks"
)

type Stdout struct {
	fd *os.File
}

func FromConfig(cfg config.Configs) sinks.Sink {
	return &Stdout{}
}

func (s *Stdout) Connect() {
	s.fd = os.Stdout
}

func (s *Stdout) Init() {}

func (s *Stdout) Send(data []byte) {
	_, _ = s.fd.Write(data)
	_, _ = s.fd.Write([]byte("\n"))
}

var _ sinks.Sink = &Stdout{}
