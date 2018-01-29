package enigma

import (
	"fmt"
	"log"
)

var reflectors = map[string]Reflector{
	"A":      MakeReflectorOrDie("EJMZALYXVBWFCRQUONTSPIKHGD"),
	"B_wide": MakeReflectorOrDie("YRUHQSLDPXNGOKMIEBFZCWVJAT"),
}

// MakeReflector turns a compact string representation of a reflector's internal
// wiring into an actual Reflector. In the string representation, position 0
// represents 'A', and its value represents the letter that 'A' connects to.
// Position 1 represents 'B', and so forth.
func MakeReflector(s string) (*Reflector, error) {
	var r Reflector
	if len(s) != len(r.mapping) {
		return nil, fmt.Errorf(
			"could not create reflector: input %v is not length %v but length %v",
			s, len(r.mapping), len(s))
	}
	for i := 0; i < len(s); i++ {
		r.mapping[i] = s[i] - 'A'
	}
	if err := ValidateReflector(r); err != nil {
		return nil, err
	}
	return &r, nil
}

// MakeReflectorOrDie does the same as MakeReflector, but instead of returning
// errors will kill the process in case of trouble.
func MakeReflectorOrDie(s string) Reflector {
	r, err := MakeReflector(s)
	if err != nil {
		log.Fatal(err)
	}
	return *r
}

// ValidateReflector returns `nil` if the given Reflector is valid, or an error
// otherwise.
func ValidateReflector(r Reflector) error {
	for i := 0; i < len(r.mapping); i++ {
		if r.mapping[i] < 0 || int(r.mapping[i]) > len(r.mapping) {
			return fmt.Errorf(
				"invalid reflector %v: position %v has invalid value %v (letter %q)",
				r.mapping, i, r.mapping[i], byte(r.mapping[i]+'A'))
		}
		to := r.mapping[i]
		if int(to) == i {
			return fmt.Errorf(
				"invalid reflector %v: position %v (letter-position %q) maps to itself",
				r.mapping, i, byte(i+'A'))
		}
		if int(r.mapping[to]) != i {
			return fmt.Errorf(
				"invalid reflector %v: %q maps to %q, but %q maps to %q",
				r.mapping, byte(i+'A'), byte(to+'A'), byte(to+'A'),
				byte(r.mapping[to]+'A'))
		}
	}
	return nil
}
