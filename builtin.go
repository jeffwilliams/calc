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
	RegisterBuiltin("binom", binom)
	RegisterBuiltin("choose", binom)
	RegisterBuiltin("bit", bit)
	RegisterBuiltin("now", now)
}
