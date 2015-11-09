package enigma

//import "fmt"
// import "log"

type logical struct {
	// The rotors in this machine, left-to-right.
	rotor     []logicalRotor
	reflector Reflector
}

type logicalRotor struct {
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

func setUpRotor(base Rotor, r *logicalRotor) {
	r.rlMapping = base.rlMapping
	for i := uint8(0); i < numLetters; i++ {
		r.lrMapping[r.rlMapping[i]] = byte(i)
	}
	r.turnoverPoints = base.turnoverPoints
}

func (e *logical) InstallRotors(rotors []Rotor) {
	e.rotor = make([]logicalRotor, len(rotors))
	for i, rotor := range rotors {
		setUpRotor(rotor, &e.rotor[i])
	}
}

func (e *logical) SetRingSettings(settings []byte) {
	for i, pos := range settings {
		e.rotor[i].ringsetting = pos - 'A'
	}
}

func (e *logical) SetRotorPositions(positions []byte) {
	for i, pos := range positions {
		e.rotor[i].rotation = pos - 'A'
	}
}

func (e *logical) InstallReflector(reflector Reflector) {
	e.reflector = reflector
}

func (e *logical) rotate() {
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

func (e *logical) KeyPress(k byte) byte {
	// Rotate the rotors for the next key press.
	e.rotate()

	// Determine the input on the stator.
	contact := k - 'A'

	// Pass through rotors, right to left.
	for i := len(e.rotor) - 1; i >= 0; i-- {
		// Connect to the next rotor.
		r := &e.rotor[i]
		contact = addRotation(r.rotation, r.ringsetting, contact)

		// Perform the mapping.
		contact = r.rlMapping[contact]

		// Connect back to the chassis.
		contact = removeRotation(r.rotation, r.ringsetting, contact)
	}

	// Pass through reflector.
	contact = e.reflector.mapping[contact]

	// Pass through rotors, left to right.
	for i := 0; i < len(e.rotor); i++ {
		// Connect to the next rotor.
		r := &e.rotor[i]
		contact = addRotation(r.rotation, r.ringsetting, contact)

		// Perform the mapping.
		contact = r.lrMapping[contact]

		// Connect back to the chassis.
		contact = removeRotation(r.rotation, r.ringsetting, contact)
	}

	return contact + 'A'
}

func NewLogical() Enigma {
	enigma := &logical{}
	return enigma
}
