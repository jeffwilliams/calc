package vmimpl

import "math/big"

func cloneInt(i *big.Int) *big.Int {
	return big.NewInt(0).Set(i)
}
