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
	// ResolvedToFnParmIndex means the name was resolved to the enclosing function
	ResolvedToFnParmIndex
	// ResolvedToAncestorFnParmIndex means the name was resolved to some function that is an ancestor of the enclosing function
	ResolvedToAncestorFnParmIndex
	ResolvedToVarNode
	ResolvedToVarSymbol
	ResolvedToFnNode
	ResolvedToFnSymbol
	ResolvedToBuiltinIndex
)

type fnAndParamIndex struct {
	node *ast.FuncDef
	parm int
}

func (c *compiler) resolveName(v ast.Parenter, name string) (typ ResolvedNameType, val interface{}, found bool) {
	found = true

	// 1. function parameters
	var node ast.Parenter
	if v != nil {
		node = v.GetParent()
	}
	height := 0

	for ; node != nil; node = node.GetParent() {
		if funcDef, ok := node.(*ast.FuncDef); ok {
			for i, arg := range funcDef.Args {
				if arg == name {
					if height == 0 {
						typ = ResolvedToFnParmIndex
						val = i
					} else {
						typ = ResolvedToAncestorFnParmIndex
						val = fnAndParamIndex{funcDef, i}
					}
					return
				}
			}
			height++
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
