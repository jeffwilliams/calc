package main

import (
	"fmt"

	"github.com/jeffwilliams/calc/ast"
	"github.com/jeffwilliams/calc/compiler"
)

func compile(moduleId string, prog string, builtinIndexes map[string]int, ref *compiler.Shared) (code *compiler.Compiled, err error) {
	var parsed interface{}
	parsed, err = Parse("last line", []byte(prog))
	if err != nil {
		err = fmt.Errorf("Parsing failed: %v", err)
		return
	}

	if _, ok := parsed.(*ast.Stmts); !ok {
		err = fmt.Errorf("Parsing result is not an AST")
		return
	}

	code, err = compiler.Compile(moduleId, parsed, builtinIndexes, nil)
	if err != nil {
		err = fmt.Errorf("Compilation failed: %v\n", err)
		return
	}
	return
}
