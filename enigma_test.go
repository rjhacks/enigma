package enigma

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MakeExampleEnigma(t *testing.T) Enigma {
	enigma := New()
	enigma.InstallRotors([]Rotor{Rotors["I"], Rotors["II"], Rotors["III"]})
	enigma.SetRingSettings([]byte{'A', 'A', 'A'})
	enigma.SetRotorPositions([]byte{'A', 'A', 'A'})
	enigma.InstallReflector(Reflectors["B"])
	// No plugboard connections.
	return enigma
}

func ResetExampleEnigma(e Enigma) {
	e.SetRotorPositions([]byte{'A', 'A', 'A'})
}

func TestBasic(t *testing.T) {
	assert := assert.New(t)
	enigma := MakeExampleEnigma(t)

	// Given that we're using Rotors I, II, III (left to right), the wide B-reflector,
	// and have all ring settings in A-position, the expected output of "AAAAA" is
	// "BDZGO" (source: https://en.wikipedia.org/wiki/Enigma_rotor_details).
	input := "AAAAA"
	encrypted := Type(enigma, input)
	assert.Equal("BDZGO", encrypted, "Wikipedia disagrees with this encryption")

	// Reset the rotor positions for decryption.
	ResetExampleEnigma(enigma)
	decrypted := Type(enigma, encrypted)
	assert.Equal(input, decrypted, "Failed to reverse encryption.")
}

func TestRingSetting(t *testing.T) {
	assert := assert.New(t)
	enigma := MakeExampleEnigma(t)

	// With the Rotors I, II, III (from left to right), wide B-reflector, all ring
	// settings in B-position, and start position AAA, typing AAAAA will produce the
	// encoded sequence EWTYX. (source: https://en.wikipedia.org/wiki/Enigma_rotor_details).
	enigma.SetRingSettings([]byte{'B', 'B', 'B'})
	input := "AAAAA"
	encrypted := Type(enigma, input)
	assert.Equal("EWTYX", encrypted, "Wikipedia disagrees with this encryption")

	// Reset the rotor positions for decryption.
	ResetExampleEnigma(enigma)
	decrypted := Type(enigma, encrypted)
	assert.Equal(input, decrypted, "Failed to reverse encryption.")
}

func TestRotation(t *testing.T) {
	assert := assert.New(t)
	enigma := MakeExampleEnigma(t)
	encrypted1 := Type(enigma, "A")
	encrypted2 := Type(enigma, "A")
	assert.NotEqual(encrypted1, encrypted2, "The first rotor isn't rotating")

	ResetExampleEnigma(enigma)
	encrypted1 = Type(enigma, strings.Repeat("A", 26))
	encrypted2 = Type(enigma, strings.Repeat("A", 26))
	assert.NotEqual(encrypted1, encrypted2, "The second rotor isn't rotating")

	ResetExampleEnigma(enigma)
	encrypted1 = Type(enigma, strings.Repeat("A", 26*26))
	encrypted2 = Type(enigma, strings.Repeat("A", 26*26))
	assert.NotEqual(encrypted1, encrypted2, "The third rotor isn't rotating")
}

func TestSingleDoubleStep(t *testing.T) {
	// Test based on the example single and double step sequence from
	// https://en.wikipedia.org/wiki/Enigma_rotor_details#Normalized_Enigma_sequences
	assert := assert.New(t)
	enig := MakeExampleEnigma(t)
	e := enig.(*enigma)

	// Normal sequence.
	e.SetRotorPositions([]byte{'A', 'A', 'U'})
	e.KeyPress('A') // Could be any key press.
	assert.Equal([]byte{'A', 'A', 'V'}, e.getRotorPositions(), "The rotor positions are wrong")
	e.KeyPress('A')
	assert.Equal([]byte{'A', 'B', 'W'}, e.getRotorPositions(), "The rotor positions are wrong")
	e.KeyPress('A')
	assert.Equal([]byte{'A', 'B', 'X'}, e.getRotorPositions(), "The rotor positions are wrong")

	// Double step sequence.
	e.SetRotorPositions([]byte{'A', 'D', 'U'}) // Normal step of right rotor.
	e.KeyPress('A')                            // Right rotor (III) goes in V - notch position.
	assert.Equal([]byte{'A', 'D', 'V'}, e.getRotorPositions(), "The rotor positions are wrong")
	e.KeyPress('A') // Right rotor steps, takes middle rotor (II) one further to E - notch position.
	assert.Equal([]byte{'A', 'E', 'W'}, e.getRotorPositions(), "The rotor positions are wrong")
	e.KeyPress('A') // Normal step of right, double step of middle, normal step of left.
	assert.Equal([]byte{'B', 'F', 'X'}, e.getRotorPositions(), "The rotor positions are wrong")
	e.KeyPress('A') // Normal step of right rotor.
	assert.Equal([]byte{'B', 'F', 'Y'}, e.getRotorPositions(), "The rotor positions are wrong")
}

func TestPlugboard(t *testing.T) {
	assert := assert.New(t)
	enigma := MakeExampleEnigma(t)
	enigma.SetPlugboard(MakePlugboard([]Pair{{'A', 'B'}, {'C', 'D'}}))

	// The same test as in `TestBasic`, but the plugboard modifies both the input
	// and the output. The input "AAAAA" becomes "BBBBB", whose output "AJLCS"
	// becomes "BJLDS".
	input := "AAAAA"
	encrypted := Type(enigma, input)
	assert.Equal("BJLDS", encrypted, "The plugboard had an unexpected effect")

	// Reset the rotor positions for decryption.
	ResetExampleEnigma(enigma)
	decrypted := Type(enigma, encrypted)
	assert.Equal(input, decrypted, "Failed to reverse encryption.")
}

func TestRealMessage1(t *testing.T) {
	// Based on http://wiki.franklinheath.co.uk/index.php/Enigma/Sample_Messages
	// From Enigma Instruction Manual, 1930
	assert := assert.New(t)

	enigma := New()
	enigma.InstallReflector(Reflectors["A"])
	enigma.InstallRotors([]Rotor{Rotors["II"], Rotors["I"], Rotors["III"]})
	enigma.SetRingSettings([]byte{'X', 'M', 'V'}) // Described as positions 24, 13, 22.
	enigma.SetPlugboard(MakePlugboard([]Pair{
		{'A', 'M'}, {'F', 'I'}, {'N', 'V'}, {'P', 'S'}, {'T', 'U'}, {'W', 'Z'}}))
	enigma.SetRotorPositions([]byte{'A', 'B', 'L'}) // Described as "message key".

	encrypted :=
		"GCDSE AHUGW TQGRK VLFGX UCALX VYMIG MMNMF DXTGN VHVRM MEVOU YFZSL RHDRR XFJWC FHUHM UNZEF RDISI KBGPM YVXUZ"
	decrypted := Type(enigma, encrypted)
	assert.Equal(
		"FEIND LIQEI NFANT ERIEK OLONN EBEOB AQTET XANFA NGSUE DAUSG ANGBA ERWAL DEXEN DEDRE IKMOS TWAER TSNEU STADT",
		decrypted, "Incorrect decryption")
}

func TestRealMessage2(t *testing.T) {
	// Based on http://www.mlb.co.jp/linux/science/genigma/enigma-referat/node4.html
	assert := assert.New(t)

	enigma := New()
	enigma.InstallReflector(Reflectors["B"]) // Assumed, not explicitly stated.
	enigma.InstallRotors([]Rotor{Rotors["II"], Rotors["I"], Rotors["V"]})
	enigma.SetRingSettings([]byte{'A', 'A', 'A'})
	enigma.SetPlugboard(MakePlugboard([]Pair{{'A', 'B'}, {'I', 'R'}, {'U', 'X'}, {'K', 'P'}}))
	enigma.SetRotorPositions([]byte{'F', 'R', 'A'})

	encrypted := "PCDAONONEBCJBOGLYMEEYGSHRYUBUJHMJOQZLEX"
	decrypted := Type(enigma, encrypted)
	assert.Equal("ANBULMEGRAZGOESTINGSTRENGGEHEIMEMELDUNG", decrypted, "Incorrect decryption")
}

// TODO: test "Operation Barbarossa, 1941" from http://wiki.franklinheath.co.uk/index.php/Enigma/Sample_Messages
