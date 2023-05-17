package main

import (
	"os"
	"strings"
)

func main() {
	message := "Go is expressive, concise, clean, and efficient. " +
		"Its concurrency mechanisms make it easy to write programs" +
		"that get the most out of multicore and networked machines, " +
		"while its novel type system enables flexible and modular program construction. " +
		"Go compiles quickly to machine code yet has the convenience of garbage" +
		"collection and the power of run-time reflection."
	// change the repeat count value to write longer data
	longString := strings.Repeat(message, 1000000)
	data := []byte(longString)

	// try to write 2 files
	// change this to write an unlimited number of files
	for i := 0; i < 2; i++ {
		writeDate(data)
	}

}

func writeDate(data []byte) {
	tmpfile, err := os.CreateTemp(".", "diskfiller")
	check(err)

	defer tmpfile.Close()

	_, err = tmpfile.Write(data)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
