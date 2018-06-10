package main

import (
	"fmt"
	"math/big"
)

func upcast(a, b interface{}) (an, bn interface{}, isInt bool) {
	an = a
	bn = b
	switch at := a.(type) {
	case *big.Int:
		switch b.(type) {
		case *big.Int:
			isInt = true
		case *big.Float:
			// b is a float. convert a to a float as well.
			an = big.NewFloat(0).SetInt(at)
		}
	case *big.Float:
		switch bt := b.(type) {
		case *big.Int:
			// b is an int. convert b to a float as well.
			bn = big.NewFloat(0).SetInt(bt)
		}
	case BigIntList:
		switch b.(type) {
		case BigIntList:
			isInt = true
		case BigFloatList:
			// b is a float list. convert a to a float list as well.
			al := make(BigFloatList, len(at))
			for i, v := range at {
				al[i] = big.NewFloat(0).SetInt(v)
			}
			an = al
		}
	case BigFloatList:
		switch bt := b.(type) {
		case BigIntList:
			// b is a int list. convert b to a float list as well.
			bl := make(BigFloatList, len(bt))
			for i, v := range bt {
				bl[i] = big.NewFloat(0).SetInt(v)
			}
			bn = bl
		}

	}

	return
}

// evalBinaryOp evaluates a simple expression of two operands and an operator.
// If both operands are Ints then the result is an Int, but if one of the operands is
// a Float the result is a Float. Effectively a Float at any point in an expression
// causes the entire evaluation to be converted to a Float. Note that the evaluated
// portions up to that point may have been calculated using integer arithmetic; this
// may lead to odd behavior for division.
func evalBinaryOp(op rune, a, b interface{}) (r interface{}, err error) {

	switch op {
	case '+':
		return add(a, b)
	case '-':
		return sub(a, b)
	case '*':
		return mul(a, b)
	case '/':
		return quo(a, b)
	case '^':
		return exp(a, b)
	case '&':
		return and(a, b)
	case '|':
		return or(a, b)
	}

	return nil, fmt.Errorf("Unsupported operation %v", op)
}

func evalUnaryOp(op rune, a interface{}) (r interface{}, err error) {

	switch op {
	case '-':
		return neg(a)
	case '~':
		return not(a)
	}

	return nil, fmt.Errorf("Unsupported operation %v", op)
}
