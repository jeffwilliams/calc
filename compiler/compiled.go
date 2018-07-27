package compiler

import (
	"bytes"
	"fmt"

	"github.com/jeffwilliams/calc/vm"
	. "github.com/jeffwilliams/calc/vmimpl"
)

// Compiled is the output of compiling some code. It contains a Main
// instruction slice for code that is not in a function, and a Shared for
// all the functions and variables.
type Compiled struct {
	// Main contains top-level code that is not in functions
	Main []vm.Instruction

	Shared
}

func (c Compiled) String(m *vm.VM) string {

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "Main:\n")

	for i, instr := range c.Main {
		fmt.Fprintf(&buf, "  %d: %s\n", i, m.InstructionString(&instr))
	}

	fmt.Fprintf(&buf, "Shared:\n")
	fmt.Fprintf(&buf, "%s\n", c.Shared.String(m))

	return buf.String()
}

// Linked returns the combined Main and Functions instructions into one complete runnable
// program. The format of the program is Main code first, followed by a halt instruction,
// followed by function code.
func (c Compiled) Linked() (code []vm.Instruction, err error) {

	lm := len(c.Main)
	lf := len(c.Functions)
	code = make([]vm.Instruction, lm+1+lf)
	copy(code, c.Main)
	code[lm] = I("halt", nil)
	copy(code[lm+1:len(code)], c.Functions)

	delta := lm + 1
	getOffset := func(sym Symbol) int {
		return sym.GetOffset() + delta
	}

	// Resolve functions in function calls. Any call instructions currently use the
	// function name as the operand instead of the offset. Change it to the offset.
	for i, v := range code {
		if unr, ok := v.Operand.(Unresolved); ok {
			tbl := c.FnSymbols
			if unr.typ == SymbolTypeVar {
				tbl = c.VarSymbols
			}

			sym, ok := tbl[unr.name]

			if unr.typ == SymbolTypeFn {
				if !ok {
					err = fmt.Errorf("Code refers to unresolved function %s", unr.name)
					return
				}
				v.Operand = getOffset(sym)
			} else {
				if !ok {
					err = fmt.Errorf("Code refers to unresolved variable %s", unr.name)
					return
				}
				v.Operand = sym.GetOffset()
			}
			code[i] = v
		}
	}
	return
}

// Link links c's Shared wih o's Shared, so that
// c contains both, and then sets c's Main to o's
// main. This is basically used for a repl.
func (c *Compiled) Link(o *Compiled) *Compiled {
	if c == nil {
		return o
	}
	c.Shared.Link(&o.Shared)
	c.Main = o.Main
	return c
}

// Compile compiles the passed AST `tree` into the Compiled `code`. The returned Compiled can
// then be linked with the results of other Compilations by linking their Shareds together if so desired.
// A final resolved program can then be obtained by calling code.Linked().
//
// builtinIndexes is used to find the indexes of the builtin functions for resolving binary operators.
func Compile(tree interface{}, builtinIndexes map[string]int, ref *Shared) (code *Compiled, err error) {
	var c compiler
	c.compile(tree, builtinIndexes, ref)
	return c.compiled, c.compileError

	return
}
