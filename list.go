package main

import "math/big"

type BigIntList []*big.Int
type BigFloatList []*big.Float

func cloneInt(i *big.Int) *big.Int {
	return big.NewInt(0).Set(i)
}

func cloneFloat(i *big.Float) *big.Float {
	return big.NewFloat(0).Set(i)
}

func clonebigInt(i *big.Int) *big.Int {
	return cloneInt(i)
}

func clonebigFloat(f *big.Float) *big.Float {
	return cloneFloat(f)
}
