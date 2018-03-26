package main

import (
	"fmt"
	"math/big"
	"testing"
)

func TestBuiltinFunc1(t *testing.T) {

	min := func(a, b *big.Int) (*big.Int, error) {
		if a.Cmp(b) < 0 {
			return a, nil
		} else {
			return b, nil
		}
	}

	f := RegisterBuiltin("min", min, "")

	m, err := f.Call([]interface{}{big.NewInt(4), big.NewInt(5)})
	if err != nil {
		t.Fatalf("Calling function failed with error: %v\n", err)
	}

	i := m.(*big.Int)

	if i.Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("Calling function returned unexpected value: %v\n", i)
	}
}

// Test a function that returns an error
func TestBuiltinFuncErr(t *testing.T) {

	fail := func(a, b *big.Int) (*big.Int, error) {
		return nil, fmt.Errorf("an error")
	}

	f := RegisterBuiltin("fail", fail, "")

	_, err := f.Call([]interface{}{big.NewInt(4), big.NewInt(5)})
	if err == nil {
		t.Fatalf("Calling function was supposed to fail with error, but it didn't\n")
	}

}

func TestBuiltinFuncWrongParamType(t *testing.T) {

	fn := func(a *big.Int) (*big.Int, error) {
		return big.NewInt(5), nil
	}

	f := RegisterBuiltin("five", fn, "")

	_, err := f.Call([]interface{}{big.NewFloat(4)})
	if err == nil {
		t.Fatalf("Calling function with wrong parameters didn't fail\n")
	}
}

func TestBuiltinFuncWrongNumParams(t *testing.T) {

	fn := func(a *big.Int) (*big.Int, error) {
		return big.NewInt(5), nil
	}

	f := RegisterBuiltin("five", fn, "")

	_, err := f.Call([]interface{}{big.NewInt(4), big.NewInt(5)})
	if err == nil {
		t.Fatalf("Calling function with wrong parameters didn't fail\n")
	}
}

func TestDefinedFuncSimple(t *testing.T) {

	parms := []string{}
	body := []byte("5")
	f := RegisterDefined("five", parms, body, "")

	m, err := f.Call([]interface{}{})
	if err != nil {
		t.Fatalf("Calling function failed with error: %v\n", err)
	}

	i := m.(*big.Int)

	if i.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("Calling function returned unexpected value: %v\n", i)
	}

	delete(Funcs, "five")
}

func TestDefinedFuncTwoParams(t *testing.T) {

	parms := []string{"x", "why"}
	body := []byte("x + why")
	f := RegisterDefined("sum", parms, body, "")

	m, err := f.Call([]interface{}{big.NewInt(3), big.NewInt(4)})
	if err != nil {
		t.Fatalf("Calling function failed with error: %v\n", err)
	}

	i := m.(*big.Int)

	if i.Cmp(big.NewInt(7)) != 0 {
		t.Fatalf("Calling function returned unexpected value: %v\n", i)
	}

	delete(Funcs, "sum")
}
