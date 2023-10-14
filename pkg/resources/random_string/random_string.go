package random_string // nolint:revive

import (
	"math/rand"
)

// RandomString is a random string generation resource.
type RandomString struct {
	size    uint32
	chars   []rune
	uniques uint32
	cache   []string
}

// New returns a new random string resource.
func New(cfg Config) *RandomString {
	return &RandomString{
		size:    cfg.Size,
		chars:   cfg.Chars,
		uniques: cfg.Uniques,
	}
}

// Cache creates a new list of random strings used by the generator.  The
// list is as long as the number of unique strings requested.
func (r *RandomString) Cache() []string {
	c := make([]string, r.uniques)

	for i := range c {
		c[i] = r.randomizer()
	}

	return c
}

// Get implements the Resource interface.
func (r *RandomString) Get() string {
	if r.uniques > 0 {
		if r.cache == nil {
			r.cache = r.Cache()
		}

		return r.cache[rand.Intn(len(r.cache))] //nolint:gosec
	}

	return r.randomizer()
}

// randomizer returns a random string of the configured size.
func (r *RandomString) randomizer() string {
	c := make([]rune, r.size)
	for i := range c {
		c[i] = r.chars[rand.Intn(len(r.chars))] //nolint:gosec
	}

	return string(c)
}
