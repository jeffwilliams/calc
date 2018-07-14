package compiler

import (
	"fmt"

	"github.com/jeffwilliams/calc/vm"
)

// Shared represents a set of compiled functions and variables and their symbol tables.
type Shared struct {
	// Functions contains code for compiled functions
	Functions []vm.Instruction
	// FnSymbols points to the functions in Functions
	FnSymbols SymbolTable
	// VarSymbols contains offsets for variables
	VarSymbols SymbolTable
}

// AddFn adds a function to the end of the Shared.
func (s *Shared) AddFn(name string, code []vm.Instruction, numArgs int) {
	sym := &FuncSymbol{
		BasicSymbol: BasicSymbol{
			Offset: len(s.Functions),
		},
		Size:    len(code),
		NumArgs: numArgs,
	}
	s.FnSymbols[name] = sym
	s.Functions = append(s.Functions, code...)
}

// AddFn removes a function from the Shared.
func (s *Shared) RemoveFn(name string) {
	sym, ok := s.FnSymbols[name]
	if ok {
		fnSym := sym.(*FuncSymbol)
		beg := s.Functions[0:sym.GetOffset()]
		end := s.Functions[sym.GetOffset()+fnSym.Size : len(s.Functions)]
		s.Functions = append(beg, end...)
		delete(s.FnSymbols, name)
	}
}

// AddVar adds a variable to the end of the Shared
func (s *Shared) AddVar(name string) (sym Symbol) {
	var ok bool
	if sym, ok = s.VarSymbols[name]; !ok {
		off := s.VarSymbols.HighestOffset() + 1
		sym = &VarSymbol{BasicSymbol{Offset: off}}
		s.VarSymbols[name] = sym
	}
	return
}

// Link links each of the arguments into s, modifying s as it goes. This combines
// multiple Shared into a single Shared that contains all of the functions and variables of
// each. The Shared s is modified to be the combined Shared; make a copy if you don't want it
// modified.
//
// If the same function is defined multiple times, the last definition (in argument order)
// is the only one kept.
//
// If the same variable is defined multiple times, only one definition is kept. The offset is
// the first one found.
func (s *Shared) Link(more ...*Shared) {

	for _, o := range more {
		for k, v := range o.FnSymbols {
			s.RemoveFn(k)
			fnSym := v.(*FuncSymbol)
			fmt.Printf("Shared.Link: o.Functions: %v\n", o.Functions)
			fmt.Printf("Shared.Link: %d-%d\n", v.GetOffset(), v.GetOffset()+fnSym.Size)
			s.AddFn(k, o.Functions[v.GetOffset():v.GetOffset()+fnSym.Size], fnSym.NumArgs)
		}

		// if variable already exists, leave it at the old location.
		for k, _ := range o.VarSymbols {
			s.AddVar(k)
		}
	}
}
