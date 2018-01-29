package enigma

// Type will press the `msg` sequence of keys on `e`, and returns
// the sequence of lights that result.
func Type(e Enigma, msg string) string {
	buffer := make([]byte, len(msg))
	for i := 0; i < len(msg); i++ {
		// Pass through spaces without running them through Enigma; they're only
		// there for human operator readability.
		if msg[i] == ' ' {
			buffer[i] = ' '
			continue
		}
		// Any real character goes through Enigma.
		buffer[i] = e.KeyPress(msg[i])
	}
	return string(buffer)
}
