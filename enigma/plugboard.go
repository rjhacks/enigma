package enigma

import (
	"fmt"
	"log"
)

// A Plugboard is much like a Reflector, in that it maps two letters to each
// other, and if 'A' maps to 'B', 'B' must map to 'A'. However, unlike a
// Reflector, a plugboard doesn't need to map every letter. An unmapped letter
// stays the same.
type Plugboard struct {
	mapping map[byte]byte
}

// AddPlugPair creates a mapping between `left` and `right`.
func (p *Plugboard) AddPlugPair(left, right byte) error {
	if p.mapping == nil {
		p.mapping = make(map[byte]byte)
	}

	if prev, mapped := p.mapping[left]; mapped {
		return fmt.Errorf(
			"Plug %q can't be mapped to %q, it was previously mapped to %q", left, right, prev)
	}
	if prev, mapped := p.mapping[right]; mapped {
		return fmt.Errorf(
			"Plug %q can't be mapped to %q, it was previously mapped to %q", right, left, prev)
	}

	p.mapping[left] = right
	p.mapping[right] = left
	return nil
}

func (p *Plugboard) mapLetter(letter byte) byte {
	// If no plugboard is present at all, letters are never mapped.
	if p == nil {
		return letter
	}

	output, mapped := p.mapping[letter]
	if !mapped {
		return letter
	}
	return output
}

// Pair represents a pair of letters to be mapped on a plugboard.
type Pair struct {
	left, right byte
}

// MakePlugboard creates a Plugboard that has the given mappings.
func makePlugboard(pairs []Pair) Plugboard {
	var plugboard Plugboard
	for _, pair := range pairs {
		if err := plugboard.AddPlugPair(pair.left, pair.right); err != nil {
			log.Fatalf("%s", err)
		}
	}
	return plugboard
}
