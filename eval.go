package main

import (
	"fmt"
	"math/big"
)

func upcast(a, b interface{}) (an, bn interface{}, isInt bool) {
	aint, ok := a.(*big.Int)
	if ok {
		bint, ok := b.(*big.Int)
		if ok {
			an = aint
			bn = bint
			isInt = true
			return
		} else {
			// b is a float. convert a to a float as well.
			bflt := b.(*big.Float)
			aflt := big.NewFloat(0).SetInt(aint)

			an = aflt
			bn = bflt
			return
		}
	} else {
		// a is a float.
		aflt := a.(*big.Float)
		bflt, ok := b.(*big.Float)
		if ok {
			an = aflt
			bn = bflt
			return
		} else {
			// b is an int. convert b to a float as well.
			bint := b.(*big.Int)
			bflt := big.NewFloat(0).SetInt(bint)

			an = aflt
			bn = bflt
			return
		}
	}

}

// evalBinaryOp evaluates a simple expression of two operands and an operator.
// If both operands are Ints then the result is an Int, but if one of the operands is
// a Float the result is a Float. Effectively a Float at any point in an expression
// causes the entire evaluation to be converted to a Float. Note that the evaluated
// portions up to that point may have been calculated using integer arithmetic; this
// may lead to odd behavior for division.
func evalBinaryOp(op rune, a, b interface{}) (r interface{}, err error) {
	a, b, isInt := upcast(a, b)

	if isInt {
		aint := a.(*big.Int)
		bint := b.(*big.Int)

		return evalInt(op, aint, bint)
	} else {
		aflt := a.(*big.Float)
		bflt := b.(*big.Float)

		return evalFloat(op, aflt, bflt)
	}
}

func evalUnaryOp(op rune, a interface{}) (r interface{}, err error) {
	aint, ok := a.(*big.Int)
	if ok {
		return evalInt(op, aint, nil)
	} else {
		aflt := a.(*big.Float)
		return evalFloat(op, aflt, nil)
	}
}

// evalFloat evaluates the operation between two Floats. For operations that are
// supported for both Ints and Floats, it uses the shared code generated from eval.genny.
func evalFloat(op rune, a, b *big.Float) (r *big.Float, err error) {
	msg := ""

	switch op {
	case '^':
		msg = "exponentiation"
	case '&':
		msg = "bitwise and"
	case '~':
		msg = "bitwise not"
	case '|':
		msg = "bitwise or"
	}

	if msg != "" {
		err = fmt.Errorf("eval: %s is only defined for integer expressions", msg)
		return
	}

	if b != nil {
		return evalBinaryOpbigFloat(op, a, b)
	} else {
		return evalUnaryOpbigFloat(op, a)
	}
}

// evalFloat evaluates the operation between two Ints. For operations that are
// supported for both Ints and Floats, it uses the shared code generated from eval.genny.
func evalInt(op rune, a, b *big.Int) (r *big.Int, err error) {
	switch op {
	case '^':
		r = a.Exp(a, b, nil)
	case '&':
		r = a.And(a, b)
	case '~':
		r = a.Not(a)
	case '|':
		r = a.Or(a, b)
	}
	if r != nil {
		return
	}

	if b != nil {
		return evalBinaryOpbigInt(op, a, b)
	} else {
		return evalUnaryOpbigInt(op, a)
	}
}
