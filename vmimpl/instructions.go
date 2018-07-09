package vmimpl

import (
	"fmt"
	"math/big"

	"github.com/jeffwilliams/calc/vm"
)

var InvalidOperandType = fmt.Errorf("operand is not valid for instruction")
var InvalidOperandValue = fmt.Errorf("operand value is out of range")
var InvalidFunction = fmt.Errorf("no function with that index")
var InvalidStackSize = fmt.Errorf("stack size is not valid for the instruction")
var InvalidAddress = fmt.Errorf("address into data segment is invalid")
var InvalidVariableType = fmt.Errorf("variable had wrong type")

// push an immediate value onto the stack
func pushOpHandler(state *vm.State, i *vm.Instruction) error {
	state.Stack.Push(i.Operand)
	return nil
}

// handle 'pushi' instruction (indirect push). Parameter to the instruction
// is a variable index (index into the data segment)
// pushes the value of the variable
func pushIndirectOpHandler(state *vm.State, i *vm.Instruction) error {
	ptr, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}

	val, ok := state.GetData(ptr)
	if !ok {
		return InvalidAddress
	}

	state.Stack.Push(val)
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

type CallBuiltinOperand struct {
	Index, NumParms int
}

func (op CallBuiltinOperand) StringWithState(s *vm.State) string {
	return fmt.Sprintf("builtin '%s', num parms %d", s.Builtins.Name(op.Index), op.NumParms)
}

// handle 'call builtin' instruction.
func callBuiltinOpHandler(state *vm.State, i *vm.Instruction) error {
	arg, ok := i.Operand.(CallBuiltinOperand)
	if !ok {
		return InvalidOperandType
	}

	fn, err := state.Builtins.Func(arg.Index)
	if err != nil || fn == nil {
		return InvalidFunction
	}

	n := arg.NumParms

	if n > len(state.Stack) {
		return InvalidStackSize
	}

	parms := make([]interface{}, n)
	for i := 0; i < n; i++ {
		parms[i] = state.Stack.Pop()
	}
	r, err := fn.Call(parms)
	if err != nil {
		return err
	}
	state.Stack.Push(r)

	return nil
}

// handle 'call' instruction. The operand is the address to call.
// sets Ip (instruction pointer) to the address of the call - 1, and pushes return address
func callOpHandler(state *vm.State, i *vm.Instruction) error {
	arg, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}

	state.Stack.Push(state.Ip)
	state.Ip = arg - 1
	return nil
}

// handle 'calli' instruction (indirect call). Parameter to the instruction
// is a variable index (index into the data segment)
// sets Ip (instruction pointer) to the address of the call - 1, and pushes return address
func callIndirectOpHandler(state *vm.State, i *vm.Instruction) error {
	ptr, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}

	argIntf, ok := state.GetData(ptr)
	if !ok {
		return InvalidAddress
	}

	arg, ok := argIntf.(int)
	if !ok {
		return InvalidVariableType
	}

	state.Stack.Push(state.Ip)
	state.Ip = arg - 1
	return nil
}

func popOpHandler(state *vm.State, i *vm.Instruction) error {
	state.Stack.Pop()
	return nil
}

func returnOpHandler(state *vm.State, i *vm.Instruction) error {
	addr := state.Stack.Pop()
	var ok bool
	state.Ip, ok = addr.(int)
	if !ok {
		// Restore state to before instruction (for debugging)
		state.Stack.Push(addr)
		return InvalidOperandType
	}

	return nil
}

// Meant for entering function.
// Set base pointer to the start of the variables in the stack.
// variables are pushed last-to-first, so they are accessed as return address = bp, arg1 = bp-1, arg2 = bp-2...
func enterOpHandler(state *vm.State, i *vm.Instruction) error {
	state.Bp = len(state.Stack) - 1
	return nil
}

// Meant for leaving function.
// Clean up function parameters, and setup stack in preparation for return.
// Assumes bp is set, and that bp+1 is return value, bp is return address, bp-1 is arg1, bp-2 is arg2,...
// Fixes stack so that top is return address, top-1 is return value, and N arguments are removed.
func leaveOpHandler(state *vm.State, i *vm.Instruction) error {
	numArgs, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}

	if state.Bp+1 >= len(state.Stack) {
		return InvalidStackSize
	}

	returnValue := state.Stack[state.Bp+1]
	returnAddr := state.Stack[state.Bp]

	state.Stack = state.Stack[0 : state.Bp-numArgs]
	state.Stack.Push(returnValue)
	state.Stack.Push(returnAddr)
	return nil
}

// push the ith function parameter, based on the current base pointer.
// variables are pushed last-to-first, so they are accessed as arg0 = bp, arg1 = bp-1, arg2 = bp-2...
func pushParmOpHandler(state *vm.State, i *vm.Instruction) error {
	index, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}
	index = state.Bp - index
	if index >= len(state.Stack) {
		return InvalidOperandValue
	}
	state.Stack.Push(state.Stack[state.Bp-index])
	return nil
}

// swap top two elements of the stack
/*func swapOpHandler(state *vm.State, i *vm.Instruction) error {
	if len(state.Stack) < 2 {
		return InvalidStackSize
	}
	state.Stack.Swap()
	return nil
}*/

func haltOpHandler(state *vm.State, i *vm.Instruction) error {
	state.Halt = true
	return nil
}
