package enigma

import (
	"fmt"
	"log"
)

var rotors = map[string]Rotor{
	"I":   MakeRotorOrDie("EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q'),
	"II":  MakeRotorOrDie("AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E'),
	"III": MakeRotorOrDie("BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V'),
	"IV":  MakeRotorOrDie("ESOVPZJAYQUIRHXLNFTGKDCMWB", 'J'),
	"V":   MakeRotorOrDie("VZBRGITYUPSDNHLXAWMJQOFECK", 'Z'),
}

// Rotor represents the configuration of a single Enigma rotor.
type Rotor struct {
	// Every rotor has 26 contacts on both the left and the right side.
	// Each of the 26 contacts on one side is connected to exactly one
	// contact on the other side. The mapping below expresses those
	// connections.
	//
	// Each contact has an index 0-25 that identifies its position on
	// its side of the rotor. The mapping below indicates which 'right'
	// contact is connected to which 'left' contact; this is the usual
	// mapping found to describe an Enigma rotor. To convert from the
	// string-based format that mapping is normally found in, use the
	// MakeRotor() method in 'util.go'. To check that your resulting
	// rotor makes sense, use ValidateRotor().
	rlMapping [numLetters]byte

	// Every rotor has different points at which it "turns over"
	// (causes the next rotor to advance one position). This mapping
	// indicates whether a given point is such a turnover point.
	turnoverPoints [numLetters]bool
}

// Reflector represents the configuration of a single Engima reflector.
type Reflector struct {
	// The reflector, unlike a rotor, has contacts on only one side,
	// and thus maps between contacts on the same side. If 'A' maps
	// to 'B', 'B' therefore must also map to 'A'.
	mapping [numLetters]byte
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
	if err := ValidateRotor(r); err != nil {
		return nil, err
	}
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
