package asm

import (
	"bufio"
	"io"

	"github.com/littlehawk93/ihex"
)

type CSC364 struct {
}

// Assemble convert incoming CSC364 assembly code into binary machine code for the CSC364 emulator. Returns the first parser error encountered
func (me CSC364) Assemble(r io.Reader, w io.Writer) error {

	var err error
	w, err = ihex.NewFileWriter(w, 16)

	scanner := bufio.NewScanner(r)

	var line string

	for scanner.Scan() {
		line = scanner.Text()

	}

	if err != nil {
		return err
	}
	return nil
}
