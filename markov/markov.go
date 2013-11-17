//package markov implements a probabalistic markov chain for generating arbitrary text
package markov

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

// Prefix is a Markov Chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

type Chain struct {
	Chain     map[string][]string
	prefixLen int
}

func (c Chain) NewPrefix() Prefix {
	return make(Prefix, c.prefixLen)
}

func NewChain(prefixLen int) (c *Chain) {
	c = new(Chain)
	c.prefixLen = prefixLen
	c.Chain = make(map[string][]string)
	return
}

func (c *Chain) Load(r io.Reader) (err error) {
	rd := bufio.NewReader(r)
	p := c.NewPrefix()
	for {
		var s string
		if _, err = fmt.Fscan(rd, &s); err != nil {
			if err != io.EOF {
				return
			}
			break
		}
		flat := strings.Join(p, " ")
		c.Chain[flat] = append(c.Chain[flat], s)
		p.Shift(s)
	}
	return
}

func (c *Chain) Generate(w io.Writer, n int) (err error) {
	p := c.NewPrefix()
	for ; n > 0; n-- {
		choices := c.Chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		if n != 1 {
			next += " "
		}
		if _, err = w.Write([]byte(next)); err != nil {
			return
		}

		p.Shift(next)
	}
	return
}
