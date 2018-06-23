package vmimpl

import (
	"github.com/jeffwilliams/calc/vm"
	"math/big"
)

type InvalidOperandTypeError int

func (e InvalidOperandTypeError) Error() string {
	return "operand is not valid for instruction"
}

var InvalidOperandType InvalidOperandTypeError

// push an immediate value onto the stack
func pushOpHandler(state *vm.State, i *vm.Instruction) error {
	state.Stack.Push(i.Operand)
	return nil
}

// integer add. add top two values in stack, and push result
func iAddOpHandler(state *vm.State, i *vm.Instruction) error {

	op1 := state.Stack.Pop()
	op2 := state.Stack.Pop()

	i1, ok := op1.(*big.Int)
	if !ok {
		return InvalidOperandType
	}
	i2, ok := op2.(*big.Int)
	if !ok {
		return InvalidOperandType
	}

	i1 = cloneInt(i1)
	state.Stack.Push(i1.Add(i1, i2))
	return nil
}

const (
	iAddOp = iota
	pushOp
)
