package main

import (
	"math/big"
	"math/rand"
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

func roll(num, sides *big.Int) (*big.Int, error) {
	sum := int64(0)
	sd := sides.Int64()
	for i := int64(0); i < num.Int64(); i++ {
		sum += rand.Int63n(sd) + 1
	}
	return big.NewInt(sum), nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
	RegisterBuiltin("binom", binom, "binmomial coeffient of (p1, p2)")
	RegisterBuiltin("choose", binom, "p1 choose p2. Same as binom")
	RegisterBuiltin("bit", bit, "return the value of bit p2 in p1, counting from 0")
	RegisterBuiltin("now", now, "return the number of milliseconds since epoch")
	RegisterBuiltin("roll", roll, "roll p1 dice each having p2 sides and sum the outcomes")
}
