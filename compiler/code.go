// The compiler package implements a compiler for generating instruction slices for running on a vm. As input
// it takes an AST, and produces Compiled objects. Compiled objects from different compilations may be linked
// together so that they share functions. A final runnable program can be obtained by calling Compiled.Linked().
//
// The main entry point to this package is the Compile() function.
package compiler

import (
	"fmt"

	"github.com/jeffwilliams/calc/ast"
	"github.com/jeffwilliams/calc/vm"
	. "github.com/jeffwilliams/calc/vmimpl"
)

type nodeCtx struct {
}

// Compiled is the output of compiling some code. It contains a Main
// instruction slice for code that is not in a function, and a Shared for
// all the functions and variables.
type Compiled struct {
	// Main contains top-level code that is not in functions
	Main []vm.Instruction

	Shared
}

type Unresolved string

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
		if name, ok := v.Operand.(Unresolved); ok {
			sym, ok := c.FnSymbols[string(name)]
			if !ok {
				err = fmt.Errorf("Code refers to unresolved symbol %s", string(name))
				return
			}

			v.Operand = getOffset(sym)
			code[i] = v
		}
	}
	return
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

func reverseParams(node interface{}, depth int) bool {
	switch t := node.(type) {
	case *ast.BinaryExpr:
		t.X, t.Y = t.Y, t.X
	}
	return true
}

type compiler struct {
	compiled       *Compiled
	compileError   error
	builtinIndexes map[string]int
	ref            *Shared
	functions      map[string]*ast.FuncDef
}

func (c *compiler) compile(tree interface{}, builtinIndexes map[string]int, ref *Shared) {
	c.compiled = &Compiled{
		Main: make([]vm.Instruction, 0, 1000),
		Shared: Shared{
			Functions:  make([]vm.Instruction, 0, 1000),
			FnSymbols:  SymbolTable{},
			VarSymbols: SymbolTable{},
		},
	}
	c.builtinIndexes = builtinIndexes
	c.ref = ref
	c.functions = make(map[string]*ast.FuncDef)

	// Find all functions defined in this compilation unit
	ast.Walk(c.buildFunctionTable, ast.Pre, tree)

	// Reverse the parameters of binary expressions so that they are pushed on the
	// stack in the right order.
	ast.Walk(reverseParams, ast.Pre, tree)

	// First generate code for each node in the tree and store it in that node's meta
	ast.Walk(c.compileNode, ast.Post, tree)

	// Then link all the code together as one block
	ast.Walk(c.linkMainNodes, ast.Post, tree)

	ast.Walk(c.linkFuncNodes, ast.Post, tree)

	return

}

func (c *compiler) buildFunctionTable(node interface{}, depth int) bool {
	if def, ok := node.(*ast.FuncDef); ok {
		c.functions[def.Name] = def
	}
	return true
}

func (c *compiler) compileNode(node interface{}, depth int) bool {
	switch t := node.(type) {
	/*
		case *Stmts:
			for _, s := range t.Stmts {
				wk(s)
			}
		case *UnaryExpr:
			wk(t.X)
		case *FuncCall:
			for _, a := range t.Args {
				wk(a)
			}
	*/
	case *ast.FuncCall:
		c.compileFuncCall(t)
	case *ast.BinaryExpr:
		c.compileBinaryExpr(t)
	case *ast.Number:
		c.compileNumber(t)
	case *ast.FuncDef:
		c.compileFuncDef(t)
	case *ast.Ident:
		c.compileIdent(t)
	case *ast.SetStmt:
		c.compileSetStmt(t)
	}

	return c.compileError == nil
}

func (c *compiler) compileNumber(v *ast.Number) {
	code := []vm.Instruction{
		I("push", v.Value),
	}

	v.SetMeta(code)
	//instructions = append(instructions, code...)
}

func (c *compiler) compileBinaryExpr(v *ast.BinaryExpr) {
	ndx, ok := c.builtinIndexes[v.Op]
	if !ok {
		c.compileError = fmt.Errorf("No builtin found for binary operator '%s'", v.Op)
		return
	}

	code := []vm.Instruction{
		I("callb", CallBuiltinOperand{Index: ndx, NumParms: 2}),
	}

	v.SetMeta(code)
	//instructions = append(instructions, code...)
}

// TODO: Not tested yet. Working on VM instructions to support this
func (c *compiler) compileFuncDef(v *ast.FuncDef) {
	// Function defs are special because they need to be placed either at the
	// beginning of memory (in which case we start the VM with a jump to the main code)
	// or placed at the end of memory.

	// Here we need to suck up all the code from our children, and leave them empty,
	// and store it on ourself. When the link happens, we'll be placed at the end
	// and added to a symbol table.
	collected := []vm.Instruction{I("enter", nil)}

	collect := func(node interface{}, depth int) bool {
		if depth == 0 {
			return true
		}

		meta := node.(ast.Metaer).GetMeta()
		if meta == nil {
			return true
		}
		code := meta.([]vm.Instruction)
		if code != nil {
			fmt.Printf("compiler.compileFuncDef.collect: depth %d appending %v\n", depth, code)
			for i, instr := range code {
				fmt.Printf("%d: %s %v\n", i, InstructTable.Name(instr.Opcode), instr.Operand)
			}

			collected = append(collected, code...)
		}
		node.(ast.Metaer).SetMeta(nil)
		return true
	}

	ast.Walk(collect, ast.Post, v)

	collected = append(collected, I("leave", len(v.Args)))
	collected = append(collected, I("return", nil))

	v.SetMeta(collected)
}

func (c *compiler) compileFuncCall(v *ast.FuncCall) {
	expectedNumArgs := -1

	if fnDefNode, ok := c.functions[v.Name]; ok {
		expectedNumArgs = len(fnDefNode.Args)
	}

	if expectedNumArgs == -1 {
		if c.ref != nil {
			if sym, ok := c.ref.FnSymbols[v.Name]; ok {
				fnSym := sym.(*FuncSymbol)
				expectedNumArgs = fnSym.NumArgs
			}
		}
	}

	if expectedNumArgs == -1 {
		c.compileError = fmt.Errorf("No function defined with name %s", v.Name)
		return
	}

	if expectedNumArgs != len(v.Args) {
		c.compileError = fmt.Errorf("Function %s expects %d arguments, but it is being called with %d", v.Name, expectedNumArgs, len(v.Args))
		return
	}

	// Reverse children
	for i := len(v.Args)/2 - 1; i >= 0; i-- {
		opp := len(v.Args) - 1 - i
		v.Args[i], v.Args[opp] = v.Args[opp], v.Args[i]
	}

	// Here we store the function name instead of the function address for the call.
	// We can't store the address yet because the symbol may be in another compilation unit,
	// so it may (a) be unresolved, or (b) be in this unit but this unit is linked after another
	// so the offset will change.
	code := []vm.Instruction{
		I("call", Unresolved(v.Name)),
	}

	v.SetMeta(code)
}

func (c *compiler) compileIdent(v *ast.Ident) {
	// Figure out what this is referring to. In order of scope it is either:
	// 1. a function parameter
	// 2. a variable in the closure of the function
	// 3. a global variable

	// First, check for function parameter
	node := v.GetParent()
	for node != nil {
		if funcDef, ok := node.(*ast.FuncDef); ok {
			for i, arg := range funcDef.Args {
				if arg == v.Name {
					code := []vm.Instruction{
						I("pushparm", i),
						I("clone", nil),
					}
					v.SetMeta(code)
					return
				}
			}
			// Only reference parameters from the first function back in the stack.
			break
		}
		node = node.GetParent()
	}

	// Must be a global var.

	// Make a slot for this variable if it doesn't exist.
	sym, ok := c.compiled.VarSymbols[v.Name]
	if !ok {
		sym = c.compiled.AddVar(v.Name)
	}

	// Add code to get the value of the variable.
	code := []vm.Instruction{
		I("load", sym.GetOffset()),
		I("clone", nil),
	}

	v.SetMeta(code)
}

func (c *compiler) compileSetStmt(v *ast.SetStmt) {
	// Make a slot for this variable if it doesn't exist.
	sym, ok := c.compiled.VarSymbols[v.Name]
	if !ok {
		sym = c.compiled.AddVar(v.Name)
	}

	code := []vm.Instruction{
		I("store", sym.GetOffset()),
	}

	v.SetMeta(code)
}

func (c *compiler) linkMainNodes(node interface{}, depth int) bool {
	if _, ok := node.(*ast.FuncDef); ok {
		// Later we will put this code at the end of the program, and put the
		// offset in a function table.
		return true
	}

	meta := node.(ast.Metaer).GetMeta()
	if meta != nil {
		code := node.(ast.Metaer).GetMeta().([]vm.Instruction)
		if code != nil {
			c.compiled.Main = append(c.compiled.Main, code...)
		}
	}
	return true
}

func (c *compiler) linkFuncNodes(node interface{}, depth int) bool {
	if def, ok := node.(*ast.FuncDef); ok {
		// Put this code at the end of the program, and put the
		// offset in a symbol table.
		if code := def.GetMeta().([]vm.Instruction); code != nil {
			c.compiled.AddFn(def.Name, code, len(def.Args))
		}
		return true
	}
	return true
}
