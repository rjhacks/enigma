package enigma

const numLetters uint8 = 26

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

// Enigma is the code version of the "human" interface of a physical Enigma
// machine. Its operations represent the actions a human might perform with
// a physical Enigma, such as installing rotors, or pressing keys. Where an
// operator would perform separate actions (e.g. installing the rotor, vs
// setting the rotor positions), the interface represents those as separate
// methods.
type Enigma interface {
	// InstallReflector places the reflector on the Engima's spindle. Which
	// reflector to use was an important secret encoded in the German code books.
	InstallReflector(reflector Reflector)

	// InstallRotors places rotors on the Engima's spindle. The rotors are listed
	// left-to-right. The internal wiring scheme of each rotor, and which set of
	// rotors would be used, were important secrets encoded in the German code
	// books.
	InstallRotors(rotors []Rotor)

	// SetRingSettings determines the offset to which the rotor rings are set.
	//
	// Each rotor can rotate its internal wiring relative to its outside
	// contacts, thereby changing the position of the wiring relative to the
	// turnover points and starting position. This is the offset to which the
	// ring is set. This "ring setting" was an important secret, encoded in the
	// German code books.
	//
	// The settings are (as in real Enigma operation) expressed as a list of
	// letters (with 'A' representing a logical offset of 0), with the first
	// setting representing the offset of the leftmost ring.
	SetRingSettings(settings []byte)

	// SetRotorPositions will rotate a rotor to a given starting position. The
	// starting position of each rotor was another important secret encoded in
	// the German code books.
	//
	// The settings are (as in real Enigma operation) expressed as a list of
	// letters (with 'A' representing a logical rotation  of 0), with the first
	// position representing the rotation of the leftmost ring.
	SetRotorPositions(positions []byte)

	// KeyPress takes the value of the key pressed on the keyboard, and returns
	// the value of the light that would light up in response.
	KeyPress(k byte) byte
}

type enigma struct {
	// In a physical Enigma's spindle (the component containing the rotors and
	// reflector), electrical signals enter from the right, pass through rotors
	// right-to-left, then through the reflector, then left-to-right through the
	// rotors again. The order of these components matters.

	// The reflector is the leftmost component in the Engima's physical spindle.
	reflector Reflector

	// The rotors in this machine, left-to-right.
	rotor []rotorConfig
}

type rotorConfig struct {
	Rotor

	// A rotor needs to map its contacts both ways, since contacts get
	// activated both left-to-right and right-to-left. The left-to-right
	// mapping is provided by the user; we compute its inverse here.
	lrMapping [numLetters]byte

	// A rotor can rotate its internal wiring relative to its outside
	// contacts, thereby changing the position of the wiring relative
	// to the turnover points and starting position. This is the offset
	// to which the ring is set.
	ringsetting uint8

	// A rotor can be in 'numLetters' different positions. We number
	// these 0..('numLetters'-1).
	rotation uint8
}

func setUpRotor(base Rotor, r *rotorConfig) {
	r.turnoverPoints = base.turnoverPoints
	r.rlMapping = base.rlMapping

	// From the rlMapping we can compute the lrMapping. The other configuration
	// values will be provided by the user later.
	for i := uint8(0); i < numLetters; i++ {
		r.lrMapping[r.rlMapping[i]] = byte(i)
	}
}

func (e *enigma) InstallRotors(rotors []Rotor) {
	e.rotor = make([]rotorConfig, len(rotors))
	for i, rotor := range rotors {
		setUpRotor(rotor, &e.rotor[i])
	}
}

func (e *enigma) SetRingSettings(settings []byte) {
	for i, pos := range settings {
		e.rotor[i].ringsetting = pos - 'A'
	}
}

func (e *enigma) SetRotorPositions(positions []byte) {
	for i, pos := range positions {
		e.rotor[i].rotation = pos - 'A'
	}
}

func (e *enigma) InstallReflector(reflector Reflector) {
	e.reflector = reflector
}

func (e *enigma) rotate() {
	// The rightmost rotor rotates fastest.
	for i := len(e.rotor) - 1; i >= 0; i-- {
		rotation := &e.rotor[i].rotation
		// Turn over only when the current rotor position is a turnover
		// point.
		turnover := e.rotor[i].turnoverPoints[*rotation]
		*rotation = (*rotation + 1) % numLetters
		if !turnover {
			break
		}
	}
}

func addRotation(rot uint8, ringsetting uint8, contact uint8) uint8 {
	// Adds 'numLetters' to ensure we're always mod-ing a positive number.
	return (contact + rot - ringsetting + numLetters) % numLetters
}

func removeRotation(rot uint8, ringsetting uint8, contact uint8) uint8 {
	// Adds '2*numLetters' to ensure we're always mod-ing a
	// positive number.
	return (contact - rot + ringsetting + 2*numLetters) % numLetters
}

func (e *enigma) KeyPress(k byte) byte {
	// Rotate the rotors for the next key press.
	e.rotate()

	// Determine the input on the stator.
	contact := k - 'A'

	// Pass through rotors, right to left.
	for i := len(e.rotor) - 1; i >= 0; i-- {
		// Connect from the chassis to the next rotor.
		r := &e.rotor[i]
		contact = addRotation(r.rotation, r.ringsetting, contact)

		// Perform the mapping.
		contact = r.rlMapping[contact]

		// Connect back to the chassis. Note that in the real Enigma there was no
		// chassis in between rotors, but doing all operations relative to the
		// 0-rotation chassis helps us keep our code sane.
		contact = removeRotation(r.rotation, r.ringsetting, contact)
	}

	// Pass through reflector.
	contact = e.reflector.mapping[contact]

	// Pass through rotors, left to right.
	for i := 0; i < len(e.rotor); i++ {
		// Connect from the chassis to the next rotor.
		r := &e.rotor[i]
		contact = addRotation(r.rotation, r.ringsetting, contact)

		// Perform the mapping.
		contact = r.lrMapping[contact]

		// Connect back to the chassis.
		contact = removeRotation(r.rotation, r.ringsetting, contact)
	}

	return contact + 'A'
}

// New creates a new Enigma machine.
func New() Enigma {
	enigma := &enigma{}
	return enigma
}
