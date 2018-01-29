package main

import (
	"fmt"
	"strconv"

	goflag "flag"

	"github.com/golang/glog"
	"github.com/rjhacks/enigma"
	"github.com/spf13/cobra"
)

/*
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rjhacks/enigma"
	"gopkg.in/urfave/cli.v1"
*/

var debugFlag bool

var reflectorFlag string
var rotorsFlag []string
var ringSettingsFlag []string
var plugPairsFlag []string
var rotorPositionsFlag []string

func crypt(cmd *cobra.Command, args []string) {
	if debugFlag {
		goflag.Set("alsologtostderr", "true")
	}
	goflag.Parse()

	e := enigma.New()

	// Install the reflector.
	{
		r, ok := enigma.Reflectors[reflectorFlag]
		if !ok {
			glog.Fatalf(
				"Reflector '%v' does not exist; options are %v",
				reflectorFlag, enigma.ReflectorNames())
		}
		e.InstallReflector(r)
		glog.Infof("Reflector: %v", reflectorFlag)
	}

	// Install the rotors.
	if len(rotorsFlag) != 3 {
		glog.Fatalf("This Enigma needs 3 rotors, but got rotors %v", rotorsFlag)
	}
	var rotors [3]enigma.Rotor
	for i, rname := range rotorsFlag {
		r, ok := enigma.Rotors[rname]
		if !ok {
			glog.Fatalf("Rotor %v does not exist; options are %v", rname, enigma.RotorNames())
		}
		rotors[i] = r
	}
	e.InstallRotors(rotors[:])
	glog.Infof("Rotors: %v", rotorsFlag)

	// Set the ring settings.
	if len(ringSettingsFlag) != 3 {
		glog.Fatalf("This Enigma needs 3 ring settings. Got ring settings %v", ringSettingsFlag)
	}
	var ringSettings [3]byte
	for i, flag := range ringSettingsFlag {
		// First attempt to interpret `setting` as a number.
		val, err := strconv.Atoi(flag)
		if err == nil {
			if val < 1 || val > 26 {
				glog.Fatalf("Got invalid ring setting number: %v", val)
			}
			ringSettings[i] = byte(val) + 'A' - 1
			continue
		}

		// Now attempt to interpret `setting` as a single character.
		if len(flag) > 1 {
			glog.Fatalf("Got invalid ring setting character: %v", flag)
		}
		b := flag[0]
		if b < 'A' || b > 'Z' {
			glog.Fatalf("Got invalid ring setting character: %v", b)
		}
		ringSettings[i] = b
	}
	e.SetRingSettings(ringSettings[:])
	glog.Infof("Ring settings: %q, %q, %q", ringSettings[0], ringSettings[1], ringSettings[2])

	// Set the plug pairs.
	var plugboard enigma.Plugboard
	for _, flag := range plugPairsFlag {
		if len(flag) != 2 {
			glog.Fatalf("All plug pairs must be 2 letters, such as 'AB'. Got: '%v'", flag)
		}
		if err := plugboard.AddPlugPair(flag[0], flag[1]); err != nil {
			glog.Fatalf("Could not add plug pair: %s", err)
		}
	}
	e.SetPlugboard(plugboard)
	glog.Infof("Plugboard: %v", plugPairsFlag)

	// Set the message key.
	if len(rotorPositionsFlag) != 3 {
		glog.Fatalf("This Enigma needs 3 rotor positions, got %v", rotorPositionsFlag)
	}
	var positions [3]byte
	for i, flag := range rotorPositionsFlag {
		if len(flag) != 1 {
			glog.Fatalf(
				"Every rotor position should be a single character, like 'A'. Got %v", rotorPositionsFlag)
		}
		b := flag[0]
		if b < 'A' || b > 'Z' {
			glog.Fatalf("Got invalid rotor position: %q", b)
		}
		positions[i] = b
	}
	e.SetRotorPositions(positions[:])
	glog.Infof("Rotor positions: %q, %q, %q", positions[0], positions[1], positions[2])

	// Finally, type the message!
	for _, arg := range args {
		out := enigma.Type(e, arg)
		if debugFlag {
			glog.Infof("%s = %s", arg, out)
		} else {
			fmt.Printf("%s ", out)
		}
	}
	fmt.Println("")
}

func main() {

	var cmdCrypt = &cobra.Command{
		Use:   "crypt [message]",
		Short: "Encrypt or decrypt a given message",
		Long: `In an Enigma, encrypting and decrypting are the same operation, just with different 
input. Use 'crypt' and pass in the message that you want to encrypt or decrypt. Use 
flags to set things like the rotors, plugboard, and so forth.`,
		Args: cobra.MinimumNArgs(1),
		Run:  crypt,
	}
	cmdCrypt.PersistentFlags().StringVar(&reflectorFlag, "reflector", "B", fmt.Sprintf(
		"The reflector called for by the code book. Options are %v",
		enigma.ReflectorNames()),
	)
	cmdCrypt.PersistentFlags().StringSliceVar(&rotorsFlag, "rotors", []string{"I", "II", "III"}, fmt.Sprintf(
		"The 3 rotors (in left-to-right order) called for by the code book. Options are %v",
		enigma.RotorNames()),
	)
	cmdCrypt.PersistentFlags().StringSliceVar(&ringSettingsFlag, "ringSettings", []string{"A", "A", "A"},
		`The ring setting for the rotors (in left-to-right order) called for by the code book. May be 
either characters (e.g. 'A') or numbers (e.g. 1)`)
	cmdCrypt.PersistentFlags().StringSliceVar(&plugPairsFlag, "plugPairs", []string{},
		`The plug pairs for the Enigma's plugboard. For example 'AB,CD' would indicate the plugboard
connects A<->B and C<->D`)
	cmdCrypt.PersistentFlags().StringSliceVar(&rotorPositionsFlag, "positions", []string{"A", "A", "A"},
		"The position of the Enigma's rotors. Also known as the 'key'.")

	var rootCmd = &cobra.Command{
		Use:   "enigma",
		Short: "A `golang` implementation of a German Wehrmacht (Army) Enigma I, circa December 1938.",
		Long: `This implementation of the Enigma aims to be true to the Enigma I, as it was in December 
1938. See usage examples at https://github.com/rjhacks/enigma.`,
	}
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Set to `true` for debug output")
	rootCmd.AddCommand(cmdCrypt)
	rootCmd.Execute()
}
