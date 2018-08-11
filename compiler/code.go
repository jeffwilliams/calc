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

// Unresolved is used as a placeholder where an instruction is generated,
// but the symbol that should be used for the operand is not yet available
type Unresolved struct {
	name string
	typ  SymbolType
}

// Fragment is the metadata for functions
type Fragment struct {
	main   []vm.Instruction
	fnName string
	fn     []vm.Instruction
}

func reverseParams(node interface{}, depth int) bool {
	switch t := node.(type) {
	case *ast.BinaryExpr:
		t.X, t.Y = t.Y, t.X
	case *ast.List:
		reverse(t.Elements)
	}
	return true
}

type compiler struct {
	compiled       *Compiled
	moduleId       string
	compileError   error
	builtinIndexes map[string]int
	ref            *Shared
	functions      map[string]*ast.FuncDef
	vars           map[string]*ast.SetStmt
	lambdaCnt      int
}

func (c *compiler) compile(moduleId string, tree interface{}, builtinIndexes map[string]int, ref *Shared) {
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
	c.vars = make(map[string]*ast.SetStmt)
	c.moduleId = moduleId

	// Find all functions and variables defined in this compilation unit
	ast.Walk(c.buildFunctionTable, ast.Pre, tree)
	ast.Walk(c.buildVarTable, ast.Pre, tree)

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
		if def.Name != "" {
			c.functions[def.Name] = def
		}
	}
	return true
}

func (c *compiler) buildVarTable(node interface{}, depth int) bool {
	if def, ok := node.(*ast.SetStmt); ok {
		c.vars[def.Name] = def
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
	case *ast.List:
		c.compileList(t)
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

func (c *compiler) compileList(v *ast.List) {
	ndx, ok := c.builtinIndexes["]"]
	if !ok {
		c.compileError = fmt.Errorf("No builtin ']' found for list construction")
		return
	}

	code := []vm.Instruction{
		I("callb", CallBuiltinOperand{Index: ndx, NumParms: len(v.Elements)}),
	}

	v.SetMeta(code)
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
	// or placed at the end of memory (in which case we halt in the main code before the
	// function defs).

	// Here we need to suck up all the code from our children, and leave them empty,
	// and store it on ourself. When the link happens, we'll be placed at the end
	// and added to a symbol table.
	collected := []vm.Instruction{
		I("enter", nil),
		I("vldac", len(v.Args)),
	}

	collect := func(node interface{}, depth int) bool {
		if depth == 0 {
			return true
		}

		meta := node.(ast.Metaer).GetMeta()
		if meta == nil {
			return true
		}

		var code []vm.Instruction
		switch t := meta.(type) {
		case []vm.Instruction:
			code = t
		case *Fragment:
			code = t.main
		default:
			panic("unknown meta type in ast")
		}
		if code != nil {
			collected = append(collected, code...)
		}
		node.(ast.Metaer).SetMeta(nil)
		return true
	}

	ast.Walk(collect, ast.Post, v)

	collected = append(collected, I("leave", nil))
	collected = append(collected, I("return", nil))

	// A function def can also be a lambda, being an unnamed function that acts as
	// a value. In this case we need to generate the function code, that gets stored
	// in the function area, and some immediate code that pushes the offset on the stack.
	// For that reason we generate two sets of instructions: the main part and the function part
	// as two separate blocks of code in a slice.
	main := []vm.Instruction{}
	if v.Name == "" {
		v.Name = c.allocLambdaId()
		main = []vm.Instruction{
			I("push", Unresolved{v.Name, SymbolTypeFn}),
		}
	}

	v.SetMeta(&Fragment{
		main:   main,
		fnName: v.Name,
		fn:     collected,
	})

}

func (c *compiler) allocLambdaId() string {
	s := fmt.Sprintf("@%s.lambda-%d", c.moduleId, c.lambdaCnt)
	c.lambdaCnt++
	return s
}

/*
func (c *compiler) genBuiltinLambdaId(name string) string {
	s := fmt.Sprintf("@lambda-builtin-%s", name)
}
*/

func (c *compiler) genBuiltinLambda(bltnName string, ndx int) (fnName string) {
	fnName = fmt.Sprintf("@builtin-lambda-%s", bltnName)

	_, _, ok := c.resolveName(nil, fnName)
	if !ok {
		code := []vm.Instruction{
			I("enter", nil),
			// -1 here means the parameter before the first on the stack,
			// being the param count
			//I("pushparm", -1),
			// Copy the number of args and the args themselves to the end of the stack
			I("reparm", nil),
			// -1 here means get the count from the stack
			I("callb", CallBuiltinOperand{Index: ndx, NumParms: -1}),
			I("leave", nil),
			I("return", nil),
		}
		c.compiled.AddFn(fnName, code, -1)
	}

	return
}

func (c *compiler) compileFuncCall(v *ast.FuncCall) {
	typ, ident, ok := c.resolveName(v, v.Name)
	if !ok {
		c.compileError = fmt.Errorf("No function or variable defined with name %s", v.Name)
		return

	}

	var code []vm.Instruction

	// TODO: Check the expected number of args for the function to see if we are passing that amount.

	switch typ {
	case ResolvedToFnParmIndex:
		code = []vm.Instruction{
			I("push", len(v.Args)),
			I("pushparm", ident.(int)),
			I("calls", nil),
		}
	case ResolvedToVarNode:
		fallthrough
	case ResolvedToVarSymbol:
		code = []vm.Instruction{
			I("push", len(v.Args)),
			I("calli", Unresolved{v.Name, SymbolTypeVar}),
		}

	case ResolvedToFnNode:
		fallthrough
	case ResolvedToFnSymbol:
		code = []vm.Instruction{
			I("push", len(v.Args)),
			I("call", Unresolved{v.Name, SymbolTypeFn}),
		}

	case ResolvedToBuiltinIndex:
		code = []vm.Instruction{
			I("callb", CallBuiltinOperand{Index: ident.(int), NumParms: len(v.Args)}),
		}
	default:
		c.compileError = fmt.Errorf("Name resolved to unknown type %d", typ)
		return

	}
	// Reverse children
	reverse(v.Args)
	v.SetMeta(code)
}

func (c *compiler) compileIdent(v *ast.Ident) {
	// Figure out what this is referring to. In order of scope it is either:
	// 1. a function parameter
	// 2. a variable in the closure of the function
	// 3. a global variable

	typ, ident, ok := c.resolveName(v, v.Name)
	if !ok {
		c.compileError = fmt.Errorf("No identifier defined with name %s", v.Name)
		return
	}

	var main []vm.Instruction
	var fn []vm.Instruction

	switch typ {
	case ResolvedToFnParmIndex:
		main = []vm.Instruction{
			I("pushparm", ident.(int)),
			I("clone", nil),
		}
	case ResolvedToVarNode:
		fallthrough
	case ResolvedToVarSymbol:
		main = []vm.Instruction{
			I("load", Unresolved{v.Name, SymbolTypeVar}),
			I("clone", nil),
		}
	case ResolvedToFnNode:
		fallthrough
	case ResolvedToFnSymbol:
		main = []vm.Instruction{
			I("push", Unresolved{v.Name, SymbolTypeFn}),
		}

	case ResolvedToBuiltinIndex:
		// TODO: Here we need to generate some sort of function F that indirectly
		// calls the builtin using it's index. We would then push the address of F here.
		// Maybe we can reserve the first N function slots as functions that just call the
		// first N builtins. Alternately we can lookup/create a new stub when we need one and
		// refer to it here.
		//c.compileError = fmt.Errorf("Identifier resolved to builtin, but no way to push builtin address onto stack.")
		shimName := c.genBuiltinLambda(v.Name, ident.(int))
		main = []vm.Instruction{
			I("push", Unresolved{shimName, SymbolTypeFn}),
		}

	default:
		c.compileError = fmt.Errorf("Name resolved to unknown type %d", typ)
		return

	}
	//v.SetMeta(&FnMeta{main: main, fn: collected})
	v.SetMeta(&Fragment{main: main, fn: fn})
	v.SetMeta(&Fragment{
		main:   main,
		fnName: v.Name,
		fn:     fn,
	})

	//v.SetMeta(code)
}

func (c *compiler) compileSetStmt(v *ast.SetStmt) {
	// Make a slot for this variable if it doesn't exist.
	c.compiled.AddVar(v.Name)

	code := []vm.Instruction{
		I("store", Unresolved{v.Name, SymbolTypeVar}),
	}

	v.SetMeta(code)
}

func (c *compiler) linkMainNodes(node interface{}, depth int) bool {
	meta := node.(ast.Metaer).GetMeta()
	if meta != nil {
		var code []vm.Instruction
		switch t := meta.(type) {
		case []vm.Instruction:
			code = t
		case *Fragment:
			code = t.main
		default:
			panic("unknown meta type in ast")
		}

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
		fnMeta := def.GetMeta().(*Fragment)

		if fnMeta.fn != nil {
			c.compiled.AddFn(fnMeta.fnName, fnMeta.fn, len(def.Args))
		}

		return true
	}
	return true
}
