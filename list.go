package main

import "math/big"

type BigIntList []*big.Int
type BigFloatList []*big.Float

func cloneInt(i *big.Int) *big.Int {
	return big.NewInt(0).Set(i)
}
