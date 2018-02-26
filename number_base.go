package main

import (
	"fmt"
	"strings"
)

// Hex, decimal, binary
type numberBase int

const (
	hexBase numberBase = iota
	decimalBase
	binaryBase
)

func (n numberBase) String() string {
	switch n {
	case hexBase:
		return "hex"
	case decimalBase:
		return "dec"
	case binaryBase:
		return "bin"
	default:
		return "unknown"
	}
}

func (n *numberBase) Set(s string) error {
	switch {
	case strings.Contains("hex", s):
		*n = hexBase
	case strings.Contains("dec", s):
		*n = decimalBase
	case strings.Contains("bin", s):
		*n = binaryBase
	default:
		return fmt.Errorf("invalid base")
	}
	return nil
}

func (n numberBase) Type() string {
	return "numberBase"
}

func (n numberBase) format(num fmt.Formatter) string {
	switch n {
	case hexBase:
		return fmt.Sprintf("0x%x", num)
	case decimalBase:
		return fmt.Sprintf("%v", num)
	case binaryBase:
		return fmt.Sprintf("0b%b", num)
	default:
		return fmt.Sprintf("%v", num)
	}
}
