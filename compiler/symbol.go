package compiler

import (
	"bytes"
	"fmt"
	"sort"
)

type SymbolType int

const (
	SymbolTypeFn SymbolType = iota
	SymbolTypeVar
)

func (s SymbolType) String() string {
	switch s {
	case SymbolTypeFn:
		return "function"
	case SymbolTypeVar:
		return "var"
	}
	return ""
}

// Symbol represents either the offset and size of a function's code in
// an instruction slice, or the offset of a variable in a data segment.
type Symbol interface {
	GetOffset() int
	SetOffset(i int)
}

type BasicSymbol struct {
	Offset int
}

func (b BasicSymbol) GetOffset() int {
	return b.Offset
}

func (b *BasicSymbol) SetOffset(i int) {
	b.Offset = i
}

type VarSymbol struct {
	BasicSymbol
}

func (s VarSymbol) String() string {
	return fmt.Sprintf("[var o:%d]", s.GetOffset())
}

type FuncSymbol struct {
	BasicSymbol
	Size, NumArgs int
}

func (s FuncSymbol) String() string {
	return fmt.Sprintf("[fn o:%d l:%d argc: %d]", s.GetOffset(), s.Size, s.NumArgs)
}

// SymbolTable maps the names of functions to their offset and size,
// or maps names of variables to their offsets.
type SymbolTable map[string]Symbol

// HighestOffset returns the offset of the symbol in the table with the highest offset
func (s SymbolTable) HighestOffset() int {
	max := -1
	for _, v := range s {
		if v.GetOffset() > max {
			max = v.GetOffset()
		}
	}
	return max
}

// AddToOffsets increases the offsets of all symbols by `delta`. Used when linking
// together code.
func (s SymbolTable) AddToOffsets(delta int) {
	for _, v := range s {
		v.SetOffset(v.GetOffset() + delta)
	}
}

func (s SymbolTable) String() string {
	var buf bytes.Buffer

	keys := make([]string, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i += 1
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(&buf, "%s: %s", k, s[k])
	}

	return buf.String()
}

// OffsetMap returns a map of an offset to the symbol name at that location.
func (s SymbolTable) OffsetMap(delta int) map[int]string {
	m := map[int]string{}
	for k, v := range s {
		m[v.GetOffset()+delta] = k
	}
	return m
}
