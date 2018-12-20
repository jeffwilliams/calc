package vmimpl

import "math/big"

var Clone func(v interface{}) (val interface{}, ok bool)

func cloneInt(i *big.Int) *big.Int {
	return big.NewInt(0).Set(i)
}

func internalClone(v interface{}) (val interface{}, ok bool) {
	ok = true
	switch v.(type) {
	case Ref:
		val = v
	default:
		ok = false
		val = v
	}
	return
}
