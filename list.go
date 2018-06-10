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

func cloneIntList(l BigIntList) BigIntList {
	l2 := make(BigIntList, len(l))
	for i, v := range l {
		l2[i] = cloneInt(v)
	}
	return l2
}

func cloneFloatList(l BigFloatList) BigFloatList {
	l2 := make(BigFloatList, len(l))
	for i, v := range l {
		l2[i] = cloneFloat(v)
	}
	return l2
}

func clone(v interface{}) interface{} {
	switch t := v.(type) {
	case *big.Int:
		return cloneInt(t)
	case *big.Float:
		return cloneFloat(t)
	case BigIntList:
		return cloneIntList(t)
	case BigFloatList:
		return cloneFloatList(t)
	}
	return v
}
