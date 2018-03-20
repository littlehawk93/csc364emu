# LA Tech CSC 364 Emulator

*Last updated on 2018-03-20*

This emulator was originally designed as a learning experience while enrolled in CSC 364 (Computer Architecture) at Louisiana Tech university. In the class, we learned the basics of digital circuitry and started from designing simple circuits (like the full adder) to a fully functional 16 bit micro-processor with a clock, ALU, ROM, and RAM.

In addition to desigining digital circuits, we learned the basics of assembly by building an assembly language based on the design of our micro-processor. That's where this project was born. In an effort to more easily practice coding our new assembly language and understand the design of the micro-processor better, I created this emulator, which I have creatively named "csc364emu".

---
## Emulator Specs 

The csc364emu micro-processor is a 16 bit processor, meaning each instruction is comprised of 2 bytes (16 bits). It has 16 registers, each (2 bytes) in size, including 4 special registers which will be further detailed later in this document. The csc364emu micro-processor has 131,070 bytes of ROM (stores 65,565 instructions) which is accessible from memory addresses 0-65,564 inclusive. The ROM memory address 65,535 or 0xFFFF is reserved to signal that csc364emu should halt. The csc364emu micro-processor has 65,536 bytes of RAM which are accessible from memory addresses 0-65,535 inclusive. Finally, csc364emu has a primitive 16x8 screen for displaying information.

---
## Instructions

Each instruction is broken down into four discrete 4 bit values:

* O (Operation)
* D (Destination)
* A (Argument A)
* B (Argument B)

#### Operation

Specifies which operation in the ALU to execute for this instruction. There are 16 different operations supported by csc364emu which will be further detailed later on in this document.

#### Destination

Specifies the address of the register the result of the operation should be stored in.

#### Argument A

The first argument used in the operation. Depending on the operation, this could be the address of a register to use as input, or just a literal binary value.

#### Argument B

The second argument used in the operation. Some operations do not support a second argument and ignore this value. Depending on the operation, this could be the address of a register to use as input, or just a literal binary value.

Each of these values are parsed left-to-right from the instruction. Consider the following instruction:

|Instruction|Binary Value|
|:---:|:---:|
 |0x45A6|01000101 10100110|

When parsed by csc364emu. It would result in this instruction values:

O|D|A|B
:---: | :---: | :---: | :---:
0100|0101|1010|0110

This instruction is basically saying: "execute operation 4 using 10 and 6 as inputs and store the result in Register 5"

---
## Operations

As stated above, csc364emu supports 16 different operations.

### MOVE
( op code 0, abbreviation: MOV )

    MOVE <regDest> <regA>

The MOVE operation copies the value from the register address *regA* and stores it in *regDest*

### NOT
( op code 1, abbreviation: *none* ) 

    NOT <regDest> <regA>

The NOT operation executes a binary NOT on the register *regA* and stores the result in *regDest*

### AND
( op code 2, abbreviation: *none* ) 

    AND <regDest> <regA> <regB>

The AND operation executes a binary AND between the register *regA* and register *regB* and stores the result in *regDest*

### OR
( op code 3, abbreviation: *none* ) 

    OR <regDest> <regA> <regB>

The OR operation executes a binary OR between the register *regA* and register *regB* and stores the result in *regDest*

### ADD
( op code 4, abbreviation: *none* ) 

    ADD <regDest> <regA> <regB>

The ADD operation adds the values in register *regA* and register *regB* and stores the result in *regDest*

### SUB
( op code 5, abbreviation: *none* ) 

    SUB <regDest> <regA> <regB>

The SUB operation subtracts the value in register *regB* from the value in register *regA* and stores the result in *regDest*

### ADDI
( op code 6, abbreviation: *none* ) 

    ADDI <regDest> <regA> <valB>

The ADDI operation adds the literal binary value of *valB* to the value in register *regA* and stores the result in *regDest*

### SUBI
( op code 7, abbreviation: *none* ) 

    SUBI <regDest> <regA> <valB>

The SUBI operation subtracts the literal binary value of *valB* from the value in register *regA* and stores the result in *regDest*

### SET
( op code 8, abbreviation: *none* ) 

    SET <regDest> <valA> <valB>

The SET operation sets the upper 8 bits of the register *regDest* as 0s and stores the binary literal of *valA* and *valB* as the lower 8 bits. *valA* is the upper 4 bits of the literal, *valB* is the lower 4 bits.

### SETH
( op code 9, abbreviation: *none* ) 

    SETH <regDest> <valA> <valB>

The SETH operation sets the upper 8 bits of the register *regDest* as the binary literal of *valA* and *valB*. *valA* is the upper 4 bits of the literal, *valB* is the lower 4 bits. The SETH operation does not modify the lower 8 bits of *regDest*. 

### INCIZ
( op code 10, abbreviation: INCZ ) 

    INCIZ <regDest> <valA> <regB>

The INCIZ operation adds the literal binary value of *valA* to the value in register *regDest* and stores the result in *regDest* if the value in register *regB* is equal to zero. Otherwise does nothing.

### DECIN
( op code 11, abbreviation: DECN ) 

    DECIN <regDest> <valA> <regB>

The DECIN operation subtracts the literal binary value of *valA* from the value in register *regDest* and stores the result in *regDest* if the value in register *regB* is negative (most significant bit is 1). Otherwise does nothing.

### MOVEZ
( op code 12, abbreviation: MOVZ ) 

    MOVEZ <regDest> <regA> <regB>

The MOVEZ operation copies the value from the register address *regA* and stores it in *regDest* if the value in register *regB* is equal to zero. Otherwise does nothing.

### MOVEX
( op code 13, abbreviation: MOVX ) 

    MOVEX <regDest> <regA> <regB>

The MOVEX operation copies the value from the register address *regA* and stores it in *regDest* if the value in register *regB* is not equal to zero. Otherwise does nothing.

### MOVEP
( op code 14, abbreviation: MOVP ) 

    MOVEP <regDest> <regA> <regB>

The MOVEP operation copies the value from the register address *regA* and stores it in *regDest* if the value in register *regB* is positive (most significant bit is 0). Otherwise does nothing.

### MOVEN
( op code 15, abbreviation: MOVN ) 

    MOVEN <regDest> <regA> <regB>

The MOVEN operation copies the value from the register address *regA* and stores it in *regDest* if the value in register *regB* is negative (most significant bit is 1). Otherwise does nothing.

---
## Registers

The csc364emu micro-processor has 16 registers designated R0 - R15. Each register can be read and written to using the operations listed above, even special registers. Registers R0 - R5 inclusive and R7 - R12 inclusive are general purpose registers that exist for storing / retrieving values. There are four special registers:

* R6 - INPUT
* R13 - OUTPUT0
* R14 - OUTPUT1
* R15 - PROGRAM COUNTER

### INPUT
( abbrevation: IN )

The INPUT register is a special register for reading values from the csc364emu's RAM or screen. When in read mode, the lower 8 bits of the INPUT register are set to the value stored at the screen or RAM address stored in the OUTPUT1 register. The high 8 bits of the INPUT register are initially set to 0s, but are not overwritten when the INPUT register is set by RAM or screen. When in read mode, the INPUT register value is set before an instruction is executed. When in write mode, the INPUT register value is written to RAM or to the screen after an instruction is executed.

### OUTPUT0
( abbreviation: OUT0 )

The OUTPUT0 register is a special register for writing values to the csc364emu's RAM or screen. Additionally, the OUTPUT0 register controls whether csc364emu's micro-processor is in read or write mode and which device to read or write to (screen or RAM). The lower 8 bits of the OUTPUT0 register are used to store the value to write to the screen or RAM. The most significant bit in the OUTPUT0 register is a flag to toggle between read and write mode (0 is read mode, 1 is write mode). The next most significant bit in the OUTPUT0 register is a flag to toggle between screen and RAM for reading from or writing to (0 is RAM, 1 is screen).

### OUTPUT1
( abbreviation: OUT1 )

The OUTPUT1 register is a special reigster for specifying which memory address from the screen or RAM to read or write to. RAM supports all memory addresses for an unsigned 16 bit integer (0-65,535 inclusive). The screen supports the values 0-15 inclusive. Any memory addresses beyond this range will not be written to and will produce a 0 if read from.

### PROGRAM COUNTER
( abbreviation: PC )

The PROGRAM COUNTER register is special register that denotes what ROM memory address the next instruction will be read from. By manipulating the PROGRAM COUNTER register, you can create conditional logics, loops, and jumps to control program flow. All ROM memory addresses in the range 0-65,534 inclusive are accepted. Setting the PROGRAM COUNTER to the maximum unsigned 16 bit integer value (65,535) signals the csc364emu to halt, stopping all execution. The PROGRAM COUNTER register is specially designed to auto-increment its value by 1 after each instruction executes unless that instruction directly manipulated the value of the PROGRAM COUNTER register (*regDest* = R15).

---
## Programming the CSC364EMU

The csc364emu accepts a ROM file as input when loaded. This ROM file is used to initialize the memory values of the csc354emu micro-processor's ROM before executing. Currently, ROM files must be formatted in Intel's [HEX](https://en.wikipedia.org/wiki/Intel_HEX) file format, specifically the I8HEX file specification. This file specification allows for specific memory addresses to be assigned in any order and features checksum validation. The csc364emu expects the final record in the provided ROM file to be a HEX end-of-file record type. Behavior is undefined if no end-of-file record is provided and no records after the end-of-file record will be read. Please visit this related [CSC364 Assembler](https://github.com/littlehawk93/csc364asm) project page to see how to compile csc364emu assembly code into a HEX ROM file. 

### Instructions

All instructions consists of either 3 or 4 tokens, separated by whitespace. Instructions are separated by newlines, so only one instruction can exist per line. The assembler is case-insensative when parsing instructions. There are three kinds of instruction tokens written in csc364emu assembly: 

* Operations
* Register Addresses
* Literal Binary Values

#### Operations

Must be the first token in the line of code. Specifies which operation to execute. The assembler accepts the full name of the operation or its abbreviation if the operation has one. For example, the operation INCIZ can be referenced with either "INCIZ" or "INCZ". The operation NOT can only be referenced with "NOT" since it has no abbreviation.

#### Register Addresses

Register addresses denote a specific register to reference for the instruction. Register address tokens begin with the letter 'R' to designate it is a register address and not a literal binary value. After the 'R', the number designates which register to reference (0-15 inclusive). This number can be written in decimal (R11) or hexadecimal (RB). Special registers can be referenced using their names or abbreviation instead of a regular register address token. For example: the INPUT register can be referenced using "R6", "INPUT", or "IN". The OUTPUT0 register can be referenced using "R13", "RD", "OUTPUT0" or "OUT0". The OUTPUT1 register can be referenced using "R14", "RE", "OUTPUT1", or "OUT1". The PROGRAM COUNTER register can be referenced using "R15", "RF", or "PC". 

#### Literal Binary Values

Literal Binary Values are exactly what they say: a designation for a literal numeric value instead of a register address. They can be expressed in decimal or hexadecimal notation. The numeric value of the literal must be within the range of 0-15 inclusive (0-F in hexadecimal). Hexadecimal literals must begin with "0x" to signal the assembler that they are written in hexadecimal notation.

### Syntax

Lines that begin with the '#' character are comments and are ignored by the assembler. Currently, comments must exist on their own lines separate from instructions or else a syntax error will be thrown

    # This is a comment. The assembler will ignore the text on this line
    ADD R1 R2 R3
    # The following instruction is commented out and won't be compiled by the assembler
    # SUB R2 R1 R3

Empty lines are also ignored by the assembler.

    ADD R1 R2 R3

    # The empty line above will be ignored by the assembler
    SUB R2 R1 R3

Both comments and empty lines do not affect the memory location where each instruction is stored in ROM. The means when modifying the PROGRAM COUNTER register, do not take into account the lines of code that are empty or comments when calculating the memory address offset to jump by in the PROGRAM COUNTER register.

    # Initialize variables
    SET R0 0 0
    SET R1 0 10

    # Begin Loop
    ADDI R0 0 1
    SUB R2 R0 R1

    # Decrement the PROGRAM COUNTER by 2 since we ignore comments and whitespace
    DECIN PC 2 R2

### Halting Program Execution

All csc364emu programs should end with the following commands to halt execution of the emulator:

    SET R0 0xF 0xF
    SETH R0 0xF 0xF
    MOVE PC R0

This sets the value of register 0 to 0xFFFF and moves it to the PROGRAM COUNTER register to signal to halt execution. If you do not include this code at the end of your program, the csc364emu is not guaranteed to halt.
