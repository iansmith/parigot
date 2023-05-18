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
	longString := strings.Repeat(message, 1000)
	data := []byte(longString)

	// change this to write an unlimited number of files
	for i := 0; i < 10000000; i++ {
		err := writeDate(data)
		if err != nil {
			log.Fatalf("Disk is full!")
		}
	}

}

func writeDate(data []byte) error {
	tmpfile, err := os.CreateTemp(".", "diskfiller")
	if err != nil {
		return err
	}

	defer tmpfile.Close()

	_, err = tmpfile.Write(data)
	if err != nil {
		return err
	}

	return nil
}
