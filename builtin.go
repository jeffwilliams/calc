package main

import "math/big"

// builtin funcs
func binom(n, k *big.Int) (*big.Int, error) {
	b := big.NewInt(0)
	return b.Binomial(n.Int64(), k.Int64()), nil
}

func bit(n, i *big.Int) (*big.Int, error) {
	return big.NewInt(int64(n.Bit(int(i.Int64())))), nil
}

func init() {
	RegisterBuiltin("binom", binom)
	RegisterBuiltin("choose", binom)
	RegisterBuiltin("bit", bit)
}
