package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"math/big"
	"os"
)

//go:generate sh -c "$GOPATH/bin/pigeon $GOPATH/src/calc/calc.peg > gen_calc.go"
//go:generate $GOPATH/bin/genny -in eval.genny -out gen_eval.go gen "Number=big.Int,big.Float"

// eval evaluates a simple expression of two operands and an operator.
// If both operands are Ints then the result is an Int, but if one of the operands is
// a Float the result is a Float. Effectively a Float at any point in an expression
// causes the entire evaluation to be converted to a Float. Note that the evaluated
// portions up to that point may have been calculated using integer arithmetic; this
// may lead to odd behavior for division.
func eval(op rune, a, b interface{}) (r interface{}, err error) {

	aint, ok := a.(*big.Int)
	if ok {
		bint, ok := b.(*big.Int)
		if ok {
			return evalInt(op, aint, bint)
		} else {
			// b is a float. convert a to a float as well.
			bflt := b.(*big.Float)
			aflt := big.NewFloat(0).SetInt(aint)

			return evalFloat(op, aflt, bflt)
		}
	} else {
		// a is a float.
		aflt := a.(*big.Float)
		bflt, ok := b.(*big.Float)
		if ok {
			return evalFloat(op, aflt, bflt)
		} else {
			// b is an int. convert b to a float as well.
			bint := b.(*big.Int)
			bflt := big.NewFloat(0).SetInt(bint)

			return evalFloat(op, aflt, bflt)
		}
	}
}

func evalFloat(op rune, a, b *big.Float) (r *big.Float, err error) {
	if op == '^' {
		err = fmt.Errorf("eval: exponentiation is only defined for integer expressions")
		return
	}

	return evalbigFloat(op, a, b)
}

func evalInt(op rune, a, b *big.Int) (r *big.Int, err error) {
	if op == '^' {
		r = a.Exp(a, b, nil)
		return
	}
	return evalbigInt(op, a, b)
}

func main() {

	rl, err := readline.New("> ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: loading readline failed: %v\n", err)
		return
	}

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				break
			}
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		parsed, err := Parse("last line", []byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Println(parsed)
	}
}
