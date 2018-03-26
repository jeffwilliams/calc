package main

import (
	"math/big"
	"time"
)

// builtin funcs
func binom(n, k *big.Int) (*big.Int, error) {
	b := big.NewInt(0)
	return b.Binomial(n.Int64(), k.Int64()), nil
}

func bit(n, i *big.Int) (*big.Int, error) {
	return big.NewInt(int64(n.Bit(int(i.Int64())))), nil
}

func now() (*big.Int, error) {
	t := time.Now()
	return big.NewInt(int64(time.Duration(t.UnixNano()) / time.Millisecond)), nil
}

func init() {
	RegisterBuiltin("binom", binom, "binmomial coeffient of (p1, p2)")
	RegisterBuiltin("choose", binom, "p1 choose p2. Same as binom")
	RegisterBuiltin("bit", bit, "return the value of bit p1, counting from 0")
	RegisterBuiltin("now", now, "return the number of milliseconds since epoch")
}
