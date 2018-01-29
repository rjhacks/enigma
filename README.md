# Enigma

A `golang` implementation of a [German Wehrmacht (Army) Enigma
I](https://en.wikipedia.org/wiki/Enigma_machine), which was used to encrypt German radio messages
during World War 2.

The Enigma encryption was cracked by a team of Allied cryptographers at Bletchly Park, notably
including Alan Turing. The machines they built to break Enigma arguably kick-started the field of
computer science.

## Installation
```sh
go install github.com/rjhacks/enigma
```

## Usage
### Command-line interface
![gif](https://i.imgur.com/56jplmt.gif)

Or decrypt a real-life message from the Enigma I manual!
```sh
$GOPATH/bin/enigma crypt  \
  --reflector=A  \
  --rotors=II,I,III \
  --ringSettings=24,13,22 \
  --plugPairs=AM,FI,NV,PS,TU,WZ \
  --positions=A,B,L \
  GCDSE AHUGW TQGRK VLFGX UCALX VYMIG MMNMF DXTGN VHVRM MEVOU YFZSL RHDRR XFJWC FHUHM UNZEF RDISI KBGPM YVXUZ
```

### As a library
If you'd like to play with the Enigma in code, you can include it directly in your programs. See
`enigma/enigma_test.go` for examples.

## Model details

This implementation of the Enigma aims to be true to the Enigma I, as it was in December
1938. The defining characteristics of this model include:
* Three rotors (although the core code actually supports any number of rotors), chosen from a set of
  five rotors, `I` through `IV`.
* A single turnover point per rotor.
* A straight connection on the entry stator (AKA: entry wheel, Eintrittswalze, ETW). Straight means
  that `A` maps to `A`, `B` maps to `B`, and so forth.
* No "Uhr", a possible extension of the plugboard. 

## References

There is a wealth of information about the Enigma on the internet, thanks to its historic status.
Some pages that were particularly helpful in building the understanding of Engima for this
implementation were:

* [Wikipedia](https://en.wikipedia.org/wiki/Enigma_machine#Reflector) for its comprehensive overview
  and its [rotor details](https://en.wikipedia.org/wiki/Enigma_rotor_details).
* [Franklin Heath](http://wiki.franklinheath.co.uk/index.php/Enigma/Sample_Messages) for its sample
  messages.
* [The Crypto Museum](http://www.cryptomuseum.com/crypto/enigma/wiring.htm), for its wiring
  schematics and clear type-by-type details.

Despite these excellent sources, there is conflicting information about whether the rotors rotate
  _before_ or _after_ the electrical contact for a letter is made. Based on example encryptions, the
  correct answer here is that the rotors rotate _before_ the contact is made.