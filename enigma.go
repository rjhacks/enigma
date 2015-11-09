package enigma

const numLetters uint8 = 26

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

type Reflector struct {
	// The reflector, unlike a rotor, has contacts on only one side,
	// and thus maps between contacts on the same side. If 'A' maps
	// to 'B', 'B' therefore must also map to 'A'.
	mapping [numLetters]byte
}

type Enigma interface {
	// For the methods below, the listed order of rotors is
	// left-to-right.
	InstallRotors(rotors []Rotor)
	SetRingSettings(settings []byte)
	SetRotorPositions(positions []byte)
	InstallReflector(reflector Reflector)
	KeyPress(k byte) byte
}
