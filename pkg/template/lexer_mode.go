package template

type mode string

const (
	TextMode       mode = "TextMode"
	ExpressionMode mode = "ExpressionMode"
	StatementMode  mode = "StatementMode"
	CommentMode    mode = "CommentMode"
	RawMode        mode = "RawMode"
)

type LexerMode []mode

func (s *LexerMode) size() int {
	return len(*s)
}

func (s *LexerMode) Mode() mode {
	size := s.size()
	if size == 0 {
		panic("internal error: empty state")
	}

	return (*s)[size-1]
}

func (s *LexerMode) Start(m mode) {
	*s = append(*s, m)
}

// Do I handle the error here?  I may be easy to accidently end the wrong one
// opaquely.
func (s *LexerMode) End(m mode) error {
	size := s.size()
	if size == 0 {
		return NoModeSelectedError
	}

	if m != (*s)[size-1] {
		return ModeNotStartedError
	}

	index := size - 1
	*s = (*s)[:index]

	return nil
}
