package compiler

import (
	"github.com/jeffwilliams/calc/ast"
)

func reverse(args []ast.Expr) {
	for i := len(args)/2 - 1; i >= 0; i-- {
		opp := len(args) - 1 - i
		args[i], args[opp] = args[opp], args[i]
	}
}
