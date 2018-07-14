package compiler

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

type FuncSymbol struct {
	BasicSymbol
	Size, NumArgs int
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
