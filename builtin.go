package main

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"
)

// builtin funcs

/*** Operators ***/

func ExpBigInt(a, b *big.Int) (r *big.Int, err error) {
	r = a.Exp(a, b, nil)
	return
}

func AndBigInt(a, b *big.Int) (r *big.Int, err error) {
	r = a.And(a, b)
	return
}

func NotBigInt(a *big.Int) (r *big.Int, err error) {
	r = a.Not(a)
	return
}

func OrBigInt(a, b *big.Int) (r *big.Int, err error) {
	r = a.Or(a, b)
	return
}

func ExpBigFloat(a, b *big.Float) (r *big.Float, err error) {
	return nil, fmt.Errorf("exponentiation is only defined for integer expressions")
}

func AndBigFloat(a, b *big.Float) (r *big.Float, err error) {
	return nil, fmt.Errorf("the 'and' operation is only defined for integer expressions")
}

func NotBigFloat(a *big.Float) (r *big.Float, err error) {
	return nil, fmt.Errorf("the 'not' operation is only defined for integer expressions")
}

func OrBigFloat(a, b *big.Float) (r *big.Float, err error) {
	return nil, fmt.Errorf("the 'or' operation is only defined for integer expressions")
}

func (l BigIntList) Exp(a, b BigIntList) (n BigIntList, err error) {
	f := func(self, a, b *big.Int) *big.Int {
		return self.Exp(a, b, big.NewInt(0))
	}
	return l.apply(a, b, f)
}
func (l BigIntList) exp(a, b BigIntList) (n BigIntList, err error) {
	return l.Exp(a, b)
}

func (l BigIntList) And(a, b BigIntList) (n BigIntList, err error) {
	return l.apply(a, b, (*big.Int).And)
}
func (l BigIntList) and(a, b BigIntList) (n BigIntList, err error) {
	return l.And(a, b)
}

func (l BigIntList) Not(a BigIntList) (n BigIntList, err error) {
	return l.applyMonadic(a, (*big.Int).Not)
}
func (l BigIntList) not(a BigIntList) (n BigIntList, err error) {
	return l.Not(a)
}

func (l BigIntList) Or(a, b BigIntList) (n BigIntList, err error) {
	return l.apply(a, b, (*big.Int).Or)
}
func (l BigIntList) or(a, b BigIntList) (n BigIntList, err error) {
	return l.Or(a, b)
}

func (l BigFloatList) Exp(a, b BigFloatList) (n BigFloatList, err error) {
	return nil, fmt.Errorf("exponentiation is only defined for integer expressions")
}
func (l BigFloatList) exp(a, b BigFloatList) (n BigFloatList, err error) {
	return l.Exp(a, b)
}

func (l BigFloatList) And(a, b BigFloatList) (n BigFloatList, err error) {
	return nil, fmt.Errorf("the 'and' operation is only defined for integer expressions")
}
func (l BigFloatList) and(a, b BigFloatList) (n BigFloatList, err error) {
	return l.And(a, b)
}

func (l BigFloatList) Not(a BigFloatList) (n BigFloatList, err error) {
	return nil, fmt.Errorf("the 'not' operation is only defined for integer expressions")
}
func (l BigFloatList) not(a BigFloatList) (n BigFloatList, err error) {
	return l.Not(a)
}

func (l BigFloatList) Or(a, b BigFloatList) (n BigFloatList, err error) {
	return nil, fmt.Errorf("the 'or' operation is only defined for integer expressions")
}
func (l BigFloatList) or(a, b BigFloatList) (n BigFloatList, err error) {
	return l.Or(a, b)
}

/*** General ***/

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

func getBytes(i *big.Int) (l BigIntList, err error) {
	b := i.Bytes()
	l = make(BigIntList, len(b))
	for i, v := range b {
		l[i] = big.NewInt(int64(v))
	}
	return
}

/*** List functions ***/

func listLen(l interface{}) (*big.Int, error) {
	switch t := l.(type) {
	case BigIntList:
		return big.NewInt(int64(len(t))), nil
	case BigFloatList:
		return big.NewInt(int64(len(t))), nil
	}
	return nil, fmt.Errorf("Unsupported type for llen")
}

func listIndex(l interface{}, i *big.Int) (interface{}, error) {
	ll, err := listLen(l)

	if err != nil {
		return nil, err
	}

	if i.Int64() < 0 || i.Cmp(ll) >= 0 {
		return nil, fmt.Errorf("Index out of range")
	}

	ndx := int(i.Uint64())

	switch t := l.(type) {
	case BigIntList:
		return cloneInt(t[ndx]), nil
	case BigFloatList:
		return cloneFloat(t[ndx]), nil
	}

	return nil, fmt.Errorf("Unsupported type for parameter 1")
}

func listReverse(l interface{}) (l2 interface{}, err error) {

	switch l.(type) {
	case BigIntList:
		return listReversebigInt(l)
	case BigFloatList:
		return listReversebigFloat(l)
	}

	return nil, fmt.Errorf("Unsupported type for parameter 1")
}

func listRepeat(e interface{}, n *big.Int) (l interface{}, err error) {
	switch t := e.(type) {
	case *big.Int:
		return listRepeatbigInt(t, n)
	case *big.Float:
		return listRepeatbigFloat(t, n)
	}

	return nil, fmt.Errorf("Unsupported type for parameter 1")
}

func unbytes(l BigIntList) (*big.Int, error) {
	bytes := make([]byte, len(l))
	for i, v := range l {
		n := v.Uint64()
		if n > 255 {
			return nil, fmt.Errorf("Value too large for byte")
		}
		bytes[i] = byte(n)
	}
	return big.NewInt(0).SetBytes(bytes), nil
}

func listMap(l interface{}, fn Func) (interface{}, error) {
	switch t := l.(type) {
	case BigIntList:
		return listMapbigIntList(t, fn)
	case BigFloatList:
		return listMapbigFloatList(t, fn)
	}

	return nil, fmt.Errorf("Unsupported type for parameter 1")
}

func listReduce(l interface{}, fn Func, memo interface{}) (interface{}, error) {
	switch t := l.(type) {
	case BigIntList:
		m, ok := memo.(*big.Int)
		if !ok {
			return nil, fmt.Errorf("Type of initial value does not match type contained in list (list contains ints)")
		}
		return listReducebigIntList(t, fn, m)
	case BigFloatList:
		m, ok := memo.(*big.Float)
		if !ok {
			return nil, fmt.Errorf("Type of initial value does not match type contained in list (list contains floats)")
		}
		return listReducebigFloatList(t, fn, m)
	}

	return nil, fmt.Errorf("Unsupported type for parameter 1")

}

/*** End List functions ***/

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

	/*** Operators ***/
	RegisterBuiltin("+", add, "return p1 + p2")
	RegisterBuiltin("-", sub, "return p1 - p2")
	RegisterBuiltin("*", mul, "return p1 * p2")
	RegisterBuiltin("/", quo, "return p1 / p2")
	RegisterBuiltin("^", exp, "return p1 ^ p2")
	RegisterBuiltin("&", and, "return p1 & p2 (bitwise and)")
	RegisterBuiltin("|", or, "return p1 | p2 (bitwise or)")
	RegisterBuiltin("~", not, "return p1 | p2 (bitwise not)")
	RegisterBuiltin("neg", neg, "return -p1 ")

	/*** General functions ***/
	RegisterBuiltin("binom", binom, "binmomial coeffient of (p1, p2)")
	RegisterBuiltin("choose", binom, "p1 choose p2. Same as binom")
	RegisterBuiltin("bit", bit, "return the value of bit p2 in p1, counting from 0")
	RegisterBuiltin("now", now, "return the number of milliseconds since epoch")
	RegisterBuiltin("roll", roll, "roll p1 dice each having p2 sides and sum the outcomes")
	RegisterBuiltin("bytes", getBytes, "return a list of each byte composing an integer")
	/*** List functions ***/
	RegisterBuiltin("llen", listLen, "return length of a list")
	RegisterBuiltin("li", listIndex, "return element at index p2 in list p1")
	RegisterBuiltin("lrev", listReverse, "return a copy of list p1 with elements in reverse order")
	RegisterBuiltin("lrp", listRepeat, "return a list consisting of p1 repeated p2 times")
	RegisterBuiltin("unbytes", unbytes, "treat the list as a list of bytes and convert it to an integer")
	RegisterBuiltin("map", listMap, "return a new list which is the result of applying the function p2 to each element in p1")
	RegisterBuiltin("reduce", listReduce, "apply a dyadic function p2 to each element in the list p1 and an accumulator (having initial value p3), returning the final value of the accumulator")
	registerStdlibMath()
}
