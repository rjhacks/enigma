package enigma

import "errors"
import "fmt"
import "log"

func Type(e Enigma, msg string) string {
	buffer := make([]byte, len(msg))
	for i := 0; i < len(msg); i++ {
		buffer[i] = e.KeyPress(msg[i])
	}
	return string(buffer)
}

func ValidateRotor(r Rotor) error {
	var seen [numLetters]bool
	for i := 0; i < len(r.rlMapping); i++ {
		if r.rlMapping[i] < 0 || int(r.rlMapping[i]) > len(r.rlMapping) {
			return errors.New(fmt.Sprintf(
				"Invalid rotor %v: position %v has invalid value %v "+
					"(letter %q)",
				r.rlMapping, i, r.rlMapping[i], byte(r.rlMapping[i]+'A')))
		}
		seen[r.rlMapping[i]] = true
	}
	for i, present := range seen {
		if !present {
			return errors.New(fmt.Sprintf(
				"Invalid rotor %v: value %v (letter %q) is missing",
				r.rlMapping, i, byte(i+'A')))
		}
	}
	return nil
}

func ValidateMultipleRotors(r []Rotor) (err error) {
	for _, rotor := range r {
		err = ValidateRotor(rotor)
		if err != nil {
			return
		}
	}
	return
}

func MakeRotor(s string, turnoverPoint byte) Rotor {
	var r Rotor
	if len(s) != len(r.rlMapping) {
		log.Fatal(fmt.Sprintf(
			"Could not create rotor: input %v is not of length %v "+
				"but of length %v",
			s, len(r.rlMapping), len(s)))
	}
	for i := 0; i < len(s); i++ {
		r.rlMapping[i] = s[i] - 'A'
	}
	r.turnoverPoints[turnoverPoint-'A'] = true
	return r
}

func ValidateReflector(r Reflector) error {
	for i := 0; i < len(r.mapping); i++ {
		if r.mapping[i] < 0 || int(r.mapping[i]) > len(r.mapping) {
			return errors.New(fmt.Sprintf(
				"Invalid reflector %v: position %v has invalid value %v "+
					"(letter %q)",
				r.mapping, i, r.mapping[i], byte(r.mapping[i]+'A')))
		}
		to := r.mapping[i]
		if int(to) == i {
			return errors.New(fmt.Sprintf(
				"Invalid reflector %v: position %v (letter-position %q) "+
					"maps to itself.",
				r.mapping, i, byte(i+'A')))
		}
		if int(r.mapping[to]) != i {
			return errors.New(fmt.Sprintf(
				"Invalid reflector %v: %q maps to %q, but %q maps to %q",
				r.mapping, byte(i+'A'), byte(to+'A'), byte(to+'A'),
				byte(r.mapping[to]+'A')))
		}
	}
	return nil
}

func MakeReflector(s string) Reflector {
	var r Reflector
	if len(s) != len(r.mapping) {
		log.Fatalf(
			"Could not create reflector: input %v is not of length %v "+
				"but of length %v", s, len(r.mapping), len(s))
	}
	for i := 0; i < len(s); i++ {
		r.mapping[i] = s[i] - 'A'
	}
	return r
}
