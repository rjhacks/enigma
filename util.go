package enigma

import "fmt"
import "log"

// Type will press the `msg` sequence of keys on `e`, and returns
// the sequence of lights that result.
func Type(e Enigma, msg string) string {
	buffer := make([]byte, len(msg))
	for i := 0; i < len(msg); i++ {
		buffer[i] = e.KeyPress(msg[i])
	}
	return string(buffer)
}

// ValidateRotor returns `nil` if the given Rotor is valid, or an error
// otherwise.
func ValidateRotor(r Rotor) error {
	var seen [numLetters]bool
	for i := 0; i < len(r.rlMapping); i++ {
		if r.rlMapping[i] < 0 || int(r.rlMapping[i]) > len(r.rlMapping) {
			return fmt.Errorf(
				"invalid rotor %v: position %v has invalid value %v (letter %q)",
				r.rlMapping, i, r.rlMapping[i], byte(r.rlMapping[i]+'A'))
		}
		seen[r.rlMapping[i]] = true
	}
	for i, present := range seen {
		if !present {
			return fmt.Errorf(
				"invalid rotor %v: value %v (letter %q) is missing",
				r.rlMapping, i, byte(i+'A'))
		}
	}
	return nil
}

// ValidateRotors returns `nil` if all given rotors are valid, or an error
// otherwise.
func ValidateRotors(r []Rotor) error {
	for _, rotor := range r {
		if err := ValidateRotor(rotor); err != nil {
			return err
		}
	}
	return nil
}

// MakeRotor turns a compact string representation of a rotor's internal wiring
// into an actual Rotor. In the string representation, position 0 represents
// 'A', and its value represents the letter that 'A' connects to. Position 1
// represents 'B', and so forth.
func MakeRotor(s string, turnoverPoint byte) (*Rotor, error) {
	var r Rotor
	if len(s) != len(r.rlMapping) {
		return nil, fmt.Errorf(
			"could not create rotor: input %v is not of length %v but of length %v",
			s, len(r.rlMapping), len(s))
	}
	for i := 0; i < len(s); i++ {
		r.rlMapping[i] = s[i] - 'A'
	}
	r.turnoverPoints[turnoverPoint-'A'] = true
	return &r, nil
}

// MakeRotorOrDie does the same as MakeRotor, but instead of returning errors
// will kill the process in case of trouble.
func MakeRotorOrDie(s string, turnoverPoint byte) Rotor {
	r, err := MakeRotor(s, turnoverPoint)
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
