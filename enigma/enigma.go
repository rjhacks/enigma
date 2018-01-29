package enigma

const numLetters uint8 = 26

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

	// SetPlugboard configures the Enigma to use the given plugboard
	// configuration. The plugboard configuration was another important secret
	// encoded in the German code books.
	SetPlugboard(plugboard Plugboard)

	// KeyPress takes the value of the key pressed on the keyboard, and returns
	// the value of the light that would light up in response.
	KeyPress(k byte) byte
}

type enigma struct {
	// The Enigma's plugboard, if any. If no plugboard is present this is nil.
	plugboard *Plugboard

	// In a physical Enigma's spindle (the component containing the rotors and
	// reflector), electrical signals enter from the right, pass through rotors
	// right-to-left, then through the reflector, then left-to-right through the
	// rotors again. The order of these components matters.
	//
	// The reflector is the leftmost component in the Engima's physical spindle.
	reflector Reflector

	// The rotors in this machine, left-to-right.
	rotor []rotorState
}

type rotorState struct {
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

func setUpRotor(base Rotor, r *rotorState) {
	r.turnoverPoints = base.turnoverPoints
	r.rlMapping = base.rlMapping

	// From the rlMapping we can compute the lrMapping. The other configuration
	// values will be provided by the user later.
	for i := uint8(0); i < numLetters; i++ {
		r.lrMapping[r.rlMapping[i]] = byte(i)
	}
}

func (e *enigma) InstallRotors(rotors []Rotor) {
	e.rotor = make([]rotorState, len(rotors))
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

func (e *enigma) getRotorPositions() []byte {
	positions := make([]byte, len(e.rotor))
	for i, rotor := range e.rotor {
		positions[i] = rotor.rotation + 'A'
	}
	return positions
}

func (e *enigma) InstallReflector(reflector Reflector) {
	e.reflector = reflector
}

func (e *enigma) SetPlugboard(plugboard Plugboard) {
	e.plugboard = &plugboard
}

func (e *enigma) rotate() {
	for i := 0; i < len(e.rotor); i++ {
		// A rotor turns when any one of the following is true:
		// - It is the rightmost rotor (which always turns).
		turn := i == len(e.rotor)-1
		// - It is in a notched position itself, and there's a rotor to its left for
		//   it to push. This condition causes the "double step" effect for (only)
		//   the middle rotor in a 3-rotor machine.
		turn = turn || (i > 0 && i < len(e.rotor)-1 && e.rotor[i].turnoverPoints[e.rotor[i].rotation])
		// - Its right neighbour is in a notched position and will push it.
		turn = turn || e.rotor[i+1].turnoverPoints[e.rotor[i+1].rotation]
		if turn {
			e.rotor[i].rotation = (e.rotor[i].rotation + 1) % numLetters
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

func (e *enigma) KeyPress(letter byte) byte {
	// Rotate the rotors for the next key press.
	e.rotate()

	// Run the key press through the plugboard.
	letter = e.plugboard.mapLetter(letter)

	// Determine the input on the stator. Before the stator, while in the keyboard/plugboard/chassis
	// it's easy to talk about each contact/wire as representing a single letter. In the rotors and
	// reflector this is harder, because the ring setting can rotate the letter-markings on the rotor
	// relative to the internal wiring. It's easier to talk about "contacts" 0-25 while we're in the
	// rotors and reflector. The stator is the conversion-point.
	contact := letter - 'A'

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

	// Pass back through the stator.
	letter = contact + 'A'

	// Second pass through the plugboard.
	letter = e.plugboard.mapLetter(letter)

	return letter
}

// New creates a new Enigma machine.
func New() Enigma {
	enigma := &enigma{}
	return enigma
}
