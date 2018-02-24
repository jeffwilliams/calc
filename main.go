package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

//go:generate sh -c "$GOPATH/bin/pigeon $GOPATH/src/calc/calc.peg > gen_calc.go"
//go:generate $GOPATH/bin/genny -in eval.genny -out gen_eval.go gen "Number=big.Int,big.Float"

func eval(op rune, a, b interface{}) (r interface{}, err error) {

	aint, ok := a.(*big.Int)
	if ok {
		bint, ok := b.(*big.Int)
		if ok {
			return evalbigInt(op, aint, bint)
		} else {
			// b is a float. convert a to a float as well.
			bflt := b.(*big.Float)
			aflt := big.NewFloat(0).SetInt(aint)

			return evalbigFloat(op, aflt, bflt)
		}
	} else {
		// a is a float.
		aflt := a.(*big.Float)
		bflt, ok := b.(*big.Float)
		if ok {
			return evalbigFloat(op, aflt, bflt)
		} else {
			// b is an int. convert b to a float as well.
			bint := b.(*big.Int)
			bflt := big.NewFloat(0).SetInt(bint)

			return evalbigFloat(op, aflt, bflt)
		}
	}
}

func main() {
	parsed, err := Parse("canned", []byte("5 + 6"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	fmt.Println(parsed)

	in := bufio.NewScanner(os.Stdin)

	for in.Scan() {
		//parsed, err := Parse("last line", in.Bytes(), Debug(true))
		parsed, err := Parse("last line", in.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Println(parsed)
	}
}
