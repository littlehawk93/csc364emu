package asm

import (
	"bufio"
	"fmt"
	"io"

	"github.com/littlehawk93/ihex"
)

// Assemble convert incoming plaintext CSC364 assembly code into binary machine code for the CSC364 emulator. Returns the first parser error encountered or any io errors
func Assemble(r io.Reader, w *ihex.FileWriter) error {

	scanner := bufio.NewScanner(r)
	lineNumber := 0
	var line string

	for scanner.Scan() {
		line = scanner.Text()
		lineNumber++
		b, err := parseLine(line)

		if err != nil {
			return fmt.Errorf("Error on line %d: %s", lineNumber, err.Error())
		}

		if _, err = w.Write(b); err != nil {
			return err
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}
