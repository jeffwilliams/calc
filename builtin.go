package main

import (
	"math"
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

func wrapFloat64FuncWith1Arg(inFn interface{}) (outFn interface{}) {
	return func(a *big.Float) (*big.Float, error) {
		af, _ := a.Float64()
		f := inFn.(func(f float64) float64)
		return big.NewFloat(f(af)), nil
	}
}

func wrapFloat64FuncWith2Arg(inFn interface{}) (outFn interface{}) {
	return func(a, b *big.Float) (*big.Float, error) {
		af, _ := a.Float64()
		bf, _ := b.Float64()
		f := inFn.(func(f, h float64) float64)
		return big.NewFloat(f(af, bf)), nil
	}
}

func registerStdlibMath() {

	reg := func(name string, fn interface{}, help string) {
		RegisterBuiltin(name, wrapFloat64FuncWith1Arg(fn), help+". This function only has the precision of a float64.")
	}

	reg2 := func(name string, fn interface{}, help string) {
		RegisterBuiltin(name, wrapFloat64FuncWith2Arg(fn), help+". This function only has the precision of a float64.")
	}

	reg("abs", math.Abs, "absolute value")
	reg("acos", math.Acos, "arccosine")
	reg("acosh", math.Acos, "inverse hyperbolic cosine")
	reg("asin", math.Asin, "arcsine")
	reg("asinh", math.Asinh, "inverse hyperbolic sine")
	reg("atan", math.Atan, "arctangent")
	reg("atanh", math.Atanh, "inverse hyperbolic tangent")
	reg("cbrt", math.Cbrt, "cube root")
	reg("ceil", math.Ceil, "ceiling")
	reg("cos", math.Cos, "cosine")
	reg("cosh", math.Cosh, "hyperbolic cosine")
	reg("erf", math.Erf, "error function")
	reg("erfc", math.Erfc, "error function compliment")
	reg("exp", math.Exp, "calculates e^p1, the base-e exponential of p1")
	reg("exp2", math.Exp2, "calculates 2^p1, the base-2 exponential of p1")
	reg("exp10", math.Pow10, "calculates 10^p1, the base-10 exponential of p1")
	reg("floor", math.Floor, "floor")
	reg("gamma", math.Gamma, "gamma function")
	reg("j0", math.J0, "order zero bessel function of the first kind")
	reg("j1", math.J1, "order one bessel function of the first kind")
	reg("log", math.Log, "natural logarithm")
	reg("log10", math.Log10, "base-10 logarithm")
	reg("log2", math.Log2, "base-2 logarithm")
	reg("sin", math.Sin, "sine")
	reg("sinh", math.Sinh, "hyperbolic sine")
	reg("sqrt", math.Sqrt, "square root")
	reg("tan", math.Tan, "tangent")
	reg("tanh", math.Tanh, "hyperbolic tangent")
	reg("y0", math.Y0, "order zero bessel function of the second kind")
	reg("y1", math.Y1, "order one bessel function of the second kind")

	reg2("hypot", math.Hypot, "calculates sqrt(p1*p1 + p2*p2)")
}

func init() {
	rand.Seed(time.Now().UnixNano())
	RegisterBuiltin("binom", binom, "binmomial coeffient of (p1, p2)")
	RegisterBuiltin("choose", binom, "p1 choose p2. Same as binom")
	RegisterBuiltin("bit", bit, "return the value of bit p2 in p1, counting from 0")
	RegisterBuiltin("now", now, "return the number of milliseconds since epoch")
	RegisterBuiltin("roll", roll, "roll p1 dice each having p2 sides and sum the outcomes")
	registerStdlibMath()
}
