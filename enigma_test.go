package enigma

import (
	//	"fmt"
	"github.com/stretchr/testify/assert"
	//	"log"
	"strings"
	"testing"
)

var rotors = []Rotor{
	MakeRotorOrDie("EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q'), // Rotor I
	MakeRotorOrDie("AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E'), // Rotor II
	MakeRotorOrDie("BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V')} // Rotor III
var ringsettings = []byte{'A', 'A', 'A'}
var positions = []byte{'A', 'A', 'A'}
var reflector = MakeReflectorOrDie("YRUHQSLDPXNGOKMIEBFZCWVJAT") // Reflector B, wide

func MakeEnigma(t *testing.T) Enigma {
	err := ValidateRotors(rotors)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ValidateReflector(reflector)
	if err != nil {
		t.Fatal(err.Error())
	}
	enigma := New()
	enigma.InstallRotors(rotors)
	enigma.SetRingSettings(ringsettings)
	enigma.SetRotorPositions(positions)
	enigma.InstallReflector(reflector)
	return enigma
}

func ResetEnigma(e Enigma) {
	e.SetRotorPositions(positions)
}

func TestBasic(t *testing.T) {
	assert := assert.New(t)
	enigma := MakeEnigma(t)

	// Given that we're using rotors I, II, III (left to right), the wide B-reflector,
	// and have all ring settings in A-position, the expected output of "AAAAA" is
	// "BDZGO" (source: https://en.wikipedia.org/wiki/Enigma_rotor_details).
	input := "AAAAA"
	encrypted := Type(enigma, input)
	assert.Equal("BDZGO", encrypted, "Wikipedia disagrees with this encryption")

	// Reset the rotor positions for decryption.
	enigma.SetRotorPositions(positions)
	decrypted := Type(enigma, encrypted)
	assert.Equal(input, decrypted, "Failed to reverse encryption.")
}

func TestRingSetting(t *testing.T) {
	assert := assert.New(t)
	enigma := MakeEnigma(t)

	// With the rotors I, II, III (from left to right), wide B-reflector, all ring
	// settings in B-position, and start position AAA, typing AAAAA will produce the
	// encoded sequence EWTYX. (source: https://en.wikipedia.org/wiki/Enigma_rotor_details).
	enigma.SetRingSettings([]byte{'B', 'B', 'B'})
	input := "AAAAA"
	encrypted := Type(enigma, input)
	assert.Equal("EWTYX", encrypted, "Wikipedia disagrees with this encryption")

	// Reset the rotor positions for decryption.
	enigma.SetRotorPositions(positions)
	decrypted := Type(enigma, encrypted)
	assert.Equal(input, decrypted, "Failed to reverse encryption.")
}

func TestRotation(t *testing.T) {
	assert := assert.New(t)
	enigma := MakeEnigma(t)
	encrypted1 := Type(enigma, "A")
	encrypted2 := Type(enigma, "A")
	assert.NotEqual(encrypted1, encrypted2, "The first rotor isn't rotating")

	ResetEnigma(enigma)
	encrypted1 = Type(enigma, strings.Repeat("A", 26))
	encrypted2 = Type(enigma, strings.Repeat("A", 26))
	assert.NotEqual(encrypted1, encrypted2, "The second rotor isn't rotating")

	ResetEnigma(enigma)
	encrypted1 = Type(enigma, strings.Repeat("A", 26*26))
	encrypted2 = Type(enigma, strings.Repeat("A", 26*26))
	assert.NotEqual(encrypted1, encrypted2, "The third rotor isn't rotating")
}

// TODO(rjh): double stepping
