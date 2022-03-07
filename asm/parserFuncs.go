package asm

import (
	"fmt"
	"strconv"
	"strings"
)

var parserFunctionsMap = map[byte][]parseToken{
	instructionMove:  {parseRegister, parseRegister},
	instructionNot:   {parseRegister, parseRegister},
	instructionAnd:   {parseRegister, parseRegister, parseRegister},
	instructionOr:    {parseRegister, parseRegister, parseRegister},
	instructionAdd:   {parseRegister, parseRegister, parseRegister},
	instructionSub:   {parseRegister, parseRegister, parseRegister},
	instructionAddi:  {parseRegister, parseRegister, parseRegister},
	instructionSubi:  {parseRegister, parseRegister, parseRegister},
	instructionSet:   {parseRegister, parseLiteral, parseLiteral},
	instructionSeth:  {parseRegister, parseLiteral, parseLiteral},
	instructionInciz: {parseRegister, parseLiteral, parseRegister},
	instructionDecin: {parseRegister, parseLiteral, parseRegister},
	instructionMovez: {parseRegister, parseRegister, parseRegister},
	instructionMovex: {parseRegister, parseRegister, parseRegister},
	instructionMovep: {parseRegister, parseRegister, parseRegister},
	instructionMoven: {parseRegister, parseRegister, parseRegister},
}

type parseToken func(token string) (byte, error)

func parseLine(line string) ([]byte, error) {

	tokens := strings.Fields(strings.ToLower(line))

	if len(tokens) == 0 || tokens[0][0] == '#' {
		return nil, nil
	}

	return parseTokens(tokens)
}

func parseTokens(tokens []string) ([]byte, error) {

	if len(tokens) == 0 {
		return nil, fmt.Errorf("Unable to parse empty line")
	}

	instruction, err := parseInstruction(tokens[0])

	if err != nil {
		return nil, err
	}

	parseFuncs, ok := parserFunctionsMap[instruction]

	if !ok {
		return nil, fmt.Errorf("No parser functions defined for instruction: %s (%X)", strings.ToUpper(tokens[0]), instruction)
	} else if len(tokens)-1 != len(parseFuncs) {
		return nil, fmt.Errorf("Instruction %s expects %d arguments. %d arguments provided", strings.ToUpper(tokens[0]), len(parseFuncs), len(tokens)-1)
	}

	outBytes := make([]byte, InstructionSize)
	outBytes[0] = instruction
	outBytes[1] = 0

	for i := 1; i < len(tokens); i++ {
		val, err := parseFuncs[i-1](tokens[i])

		if err != nil {
			return nil, err
		}

		if i%2 == 0 {
			outBytes[i/2] = (val << 4) & 0xF0
		} else {
			outBytes[i/2] = (outBytes[i/2] & 0xF0) | (val & 0x0F)
		}
	}
	return outBytes, nil
}

func parseInstruction(token string) (byte, error) {

	if val, ok := parseMapOrLiteral(token, instructionsMap); ok {
		return val, nil
	}
	return 0, newTokenParseError("instruction", token)
}

func parseRegister(token string) (byte, error) {

	if val, ok := parseMapOrLiteral(token, registersMap); ok {
		return val, nil
	}
	return 0, newTokenParseError("register address", token)
}

func parseMapOrLiteral(token string, m map[string]byte) (byte, bool) {

	val, ok := m[token]
	if !ok {
		var err error
		if val, err = parseLiteral(token); err != nil {
			return val, false
		}
	}
	return val, true
}

func parseLiteral(token string) (byte, error) {

	val, err := parseHex(token)

	if err != nil {
		if val, err = parseDecimal(token); err != nil {
			return 0, newTokenParseError("numeric literal", token)
		}
	}
	return val, nil
}

func parseDecimal(token string) (byte, error) {

	val, err := strconv.ParseInt(token, 10, 8)

	if err != nil {
		return 0, newTokenParseError("decimal literal", token)
	}
	return byte(val), nil
}

func parseHex(token string) (byte, error) {

	if !strings.HasPrefix(token, "0x") || len(token) <= 2 {
		return 0, newTokenParseError("hex literal", token)
	}

	val, err := strconv.ParseInt(token[2:], 16, 8)

	if err != nil {
		return 0, newTokenParseError("hex literal", token)
	}
	return byte(val), nil
}

func newTokenParseError(dataType, token string) error {
	return fmt.Errorf("Unable to parse %s from token: '%s'", dataType, token)
}
