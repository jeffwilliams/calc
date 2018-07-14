package vmimpl

import "math/big"

var Clone func(v interface{}) (val interface{}, ok bool)

func cloneInt(i *big.Int) *big.Int {
	return big.NewInt(0).Set(i)
}
