package random_string

import (
	"math/rand"

	"ctx.sh/genie/pkg/config"
)

type RandomString struct {
	size    uint32
	chars   []rune
	uniques uint32
	cache   []string
}

func FromConfig(options config.RandomString) *RandomString {
	return &RandomString{
		size:    options.Size,
		chars:   options.Chars,
		uniques: options.Uniques,
	}
}

func (r *RandomString) Cache() []string {
	c := make([]string, r.uniques)

	for i := range c {
		c[i] = r.randomizer()
	}

	return c
}

func (r *RandomString) Get() string {
	if r.uniques > 0 {
		if r.cache == nil {
			r.cache = r.Cache()
		}

		return r.cache[rand.Intn(len(r.cache))]
	}

	return r.randomizer()
}

func (r *RandomString) randomizer() string {
	c := make([]rune, r.size)
	for i := range c {
		c[i] = r.chars[rand.Intn(len(r.chars))]
	}

	return string(c)
}
