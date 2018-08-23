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
var InvalidArgumentCount = fmt.Errorf("wrong number of arguments passed to function")

// push an immediate value onto the stack
func pushOpHandler(state *vm.State, i *vm.Instruction) error {
	state.Stack.Push(i.Operand)
	return nil
}

// handle 'load' instruction.  Parameter to the instruction
// is a variable index (index into the data segment)
// pushes the value of the variable
func loadOpHandler(state *vm.State, i *vm.Instruction) error {
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

// handle 'store' instruction.  Parameter to the instruction
// is a variable index (index into the data segment)
// pops the top of the stack and stores it into the variable
func storeOpHandler(state *vm.State, i *vm.Instruction) error {
	ptr, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}

	val := state.Stack.Pop()

	ok = state.SetData(ptr, val)
	if !ok {
		// restore stack
		state.Stack.Push(val)
		return InvalidAddress
	}

	return nil
}

// integer add. add top two values in stack, and push result
func iAddOpHandler(state *vm.State, i *vm.Instruction) error {

	op1 := state.Stack.Pop()
	op2 := state.Stack.Pop()

	i1, ok := op1.(*big.Int)
	if !ok {
		// restore stack
		state.Stack.Push(op2)
		state.Stack.Push(op1)
		return InvalidOperandType
	}
	i2, ok := op2.(*big.Int)
	if !ok {
		// restore stack
		state.Stack.Push(op2)
		state.Stack.Push(op1)
		return InvalidOperandType
	}

	i1 = cloneInt(i1)
	state.Stack.Push(i1.Add(i1, i2))
	return nil
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
	if n < 0 {
		// num parms not specified in operand.
		// it must be on the stack
		n = state.Stack.Pop().(int)
	}

	if n > len(state.Stack) {
		// restore stack
		if arg.NumParms < 0 {
			state.Stack.Push(n)
		}
		return InvalidStackSize
	}

	parms := make([]interface{}, n)
	for i := 0; i < n; i++ {
		parms[i] = state.Stack.Pop()
	}
	r, err := fn.Call(parms)
	if err != nil {
		// Restore stack
		for i, _ := range parms {
			state.Stack.Push(parms[len(parms)-i-1])
		}
		if arg.NumParms < 0 {
			state.Stack.Push(n)
		}
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

	var newIp int
	switch t := argIntf.(type) {
	case int:
		newIp = t
	case LambdaClosureOperand:
		newIp = t.LambdaAddr
		// Push the pointer to the closure env as the first argument to the function.
		state.Stack.Push(t.ClosureEnv)
	default:
		return InvalidVariableType
	}

	/*arg, ok := argIntf.(int)
	if !ok {
		return InvalidVariableType
	}*/

	state.Stack.Push(state.Ip)
	state.Ip = newIp - 1
	return nil
}

// handle 'calls' instruction (call from stack). Calls the function at the address in the top of
// the stack.
// sets Ip (instruction pointer) to the address of the call - 1, and pushes return address
func callStackOpHandler(state *vm.State, i *vm.Instruction) error {
	addr := state.Stack.Pop()

	var newIp int
	switch t := addr.(type) {
	case int:
		newIp = t
	case LambdaClosureOperand:
		newIp = t.LambdaAddr
		// Push the pointer to the closure env as the first argument to the function.
		state.Stack.Push(t.ClosureEnv)
	default:
		state.Stack.Push(addr)
		return InvalidVariableType
	}

	state.Stack.Push(state.Ip)
	state.Ip = newIp - 1
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
// Set base pointer to the index in the stack of the argument count (immediately below which are the variables)
// variables are pushed last-to-first, so they are accessed as return address = bp, arg0 = bp-1, arg1 = bp-2...
// After this call the frame is: [...][arg1][arg0][argcnt][return addr][old bp]
func enterOpHandler(state *vm.State, i *vm.Instruction) error {
	// Save Bp
	state.Stack.Push(state.Bp)
	state.Bp = len(state.Stack) - 3
	return nil
}

func validateArgCount(state *vm.State, i *vm.Instruction) error {
	exp, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}
	if exp < 0 {
		return nil
	}
	numArgs := state.Stack[state.Bp]
	if exp != numArgs {
		return InvalidArgumentCount
	}
	return nil
}

// Meant for leaving function.
// Clean up function parameters, and setup stack in preparation for return.
// Fixes stack so that top is return address, top-1 is return value, and N arguments are removed.
func leaveOpHandler(state *vm.State, i *vm.Instruction) error {
	if state.Bp+3 >= len(state.Stack) {
		return InvalidStackSize
	}

	returnValue := state.Stack[state.Bp+3]
	oldBp := state.Stack[state.Bp+2].(int)
	returnAddr := state.Stack[state.Bp+1]
	numArgs := state.Stack[state.Bp].(int)

	state.Stack = state.Stack[0 : state.Bp-numArgs]
	state.Stack.Push(returnValue)
	state.Stack.Push(returnAddr)
	// Restore Bp
	state.Bp = oldBp
	return nil
}

// push the ith function parameter, based on the current base pointer.
// variables are pushed last-to-first, so they are accessed as arg0 = bp-1, arg2 = bp-2...
func pushParmOpHandler(state *vm.State, i *vm.Instruction) error {
	index, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}
	index = state.Bp - index - 1
	if index >= len(state.Stack) {
		return InvalidOperandValue
	}
	state.Stack.Push(state.Stack[index])
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

func cloneOpHandler(state *vm.State, i *vm.Instruction) error {
	top := state.Stack.Pop()
	val, ok := Clone(top)
	if !ok {
		state.Stack.Push(top)
		return InvalidOperandType
	}
	state.Stack.Push(val)
	return nil
}

// Copy a segment of the stack to the end of the stack
func copyStackOpHandler(state *vm.State, i *vm.Instruction) error {
	arg, ok := i.Operand.(CopyStackOperand)
	if !ok {
		return InvalidOperandType
	}

	state.Stack.CopyToEnd(arg.Offset, arg.Len)
	return nil
}

// Push the argument count and arguments back onto the top of the stack
// as if the function had just been called.
func reparmOpHandler(state *vm.State, i *vm.Instruction) error {
	offset := len(state.Stack) - state.Bp - 1
	len := state.Stack[state.Bp].(int) + 1

	state.Stack.CopyToEnd(offset, len)
	return nil

}
