package asm

import "io"

// Assembler defines a type that assembles an assembly code into machine code
type Assembler interface {
	Assemble(in io.Reader, out io.Writer) []error
}
