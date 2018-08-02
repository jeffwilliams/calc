package compiler

import (
	"github.com/jeffwilliams/calc/ast"
)

// resolveName tries to find the item that the identifier is referencing.
// In general it finds the first matching name in order of scope. It checks in order:
//  1. function parameters
//  2. variables in the closure of the function
//  3. global variables
//  4. function names
//
// This function returns one of:
//   nil, if the name was not found
//   an *ast.FuncDef

type ResolvedNameType int

const (
	ResolvedToNone ResolvedNameType = iota
	ResolvedToFnParmIndex
	ResolvedToVarNode
	ResolvedToVarSymbol
	ResolvedToFnNode
	ResolvedToFnSymbol
	ResolvedToBuiltinIndex
)

func (c *compiler) resolveName(v ast.Parenter, name string) (typ ResolvedNameType, val interface{}, found bool) {
	found = true

	// 1. function parameters
	node := v.GetParent()
	for ; node != nil; node = node.GetParent() {
		if funcDef, ok := node.(*ast.FuncDef); ok {
			for i, arg := range funcDef.Args {
				if arg == name {
					typ = ResolvedToFnParmIndex
					val = i
					return
				}
			}
			// Only reference parameters from the first function back in the stack.
			break
		}
	}

	// 3. global variables

	// check for variable in this module
	if varNode, ok := c.vars[name]; ok {
		typ = ResolvedToVarNode
		val = varNode
		return
	}

	// check for variable in other referenced modules
	if c.ref != nil {
		if sym, ok := c.ref.VarSymbols[name]; ok {
			typ = ResolvedToVarSymbol
			val = sym
			return
		}
	}

	// 4. Function names
	// check for fn in this module
	if fnDefNode, ok := c.functions[name]; ok {
		typ = ResolvedToFnNode
		val = fnDefNode
		return
	}

	// check for fn in other modules
	if c.ref != nil {
		if sym, ok := c.ref.FnSymbols[name]; ok {
			typ = ResolvedToFnSymbol
			val = sym
			return
		}
	}

	// check for builtin
	ndx, ok := c.builtinIndexes[name]
	if ok {
		typ = ResolvedToBuiltinIndex
		val = ndx
		return
	}

	found = false
	return
}
