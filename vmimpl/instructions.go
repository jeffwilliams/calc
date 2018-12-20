package vmimpl

import (
	"fmt"
	"math/big"

	"github.com/jeffwilliams/calc/vm"
)

// Use the VM's general register 0 for storing the current closure's environment
const ClosureEnvGenReg = 0

var instructionDesc = []vm.InstructionDescr{
	{"invalid", haltOpHandler},
	{"iadd", iAddOpHandler},
	{"push", pushOpHandler},
	{"pop", popOpHandler},
	{"swap", swapOpHandler},
	{"callb", callBuiltinOpHandler},
	{"call", callOpHandler},
	{"calli", callIndirectOpHandler},
	{"calls", callStackOpHandler},
	{"return", returnOpHandler},
	{"enter", enterOpHandler},
	{"leave", leaveOpHandler},
	{"reparm", reparmOpHandler},
	{"vldac", validateArgCount},
	{"pushparm", pushParmOpHandler},
	{"copys", copyStackOpHandler},
	{"halt", haltOpHandler},
	{"clone", cloneOpHandler},
	{"load", loadOpHandler},
	{"store", storeOpHandler},
	{"stores", storeStackOpHandler},
	{"tload", tloadOpHandler},
	{"tstore", tstoreOpHandler},
	{"tmake", tmakeOpHandler},
	{"alloc", allocOpHandler},
	{"free", freeOpHandler},
	{"mkclsr", makeClosureOpHandler},
	{"setgr", setGenRegisterOpHandler},
	{"getgr", getGenRegisterOpHandler},
}

var InstructionHelp = map[string]string{
	"invalid":  "Invalid instruction",
	"iadd":     "Pop two elements from the stack, add them, and push the result. Stack elements must be type *big.Int",
	"push":     "Push the operand onto the stack",
	"pop":      "Pop the top element of the stack, and discard it",
	"swap":     "Swap the top two elements of the stack",
	"callb":    "Call the builtin function operand.Index, popping off operand.NumParms arguments from the stack. Push the result onto the stack.",
	"call":     "Call the function at the instruction index contained in the operand. The current Ip is pushed, then Ip is set to the operand.",
	"calli":    "Call a function indirectly. The operand is treated an index into the data segment, and the value in that data slot is used as the address to call. If the data slot contains an int, the current Ip is pushed, then Ip is set to Data[operand]. If the data slot contains a LambdaClosureOperand, Data[operand].ClosureEnv is pushed, then Ip is pushed, then Ip is set to Data[operand].LambdaAddr.",
	"calls":    "Call a function who's address is on the top of stack. Same as calli, except top of stack is used instead of the data segment.",
	"return":   "Return from a function call. Top of stack is popped and Ip is set to that value. Top of stack must be an int.",
	"enter":    "Enter a function (setup stack frame). Current value of Bp is pushed and Bp is set to point to the stack slot containing the function argument count, below which are the arguments.",
	"leave":    "Leave a function (cleanup stack frame). Fixes stack so that top is return address, top-1 is return value, and the function arguments are removed from the stack.",
	"reparm":   "Push the argument count and arguments back onto the top of the stack as if the function had just been called",
	"vldac":    "Validate argument count. If the operand does not equal the value at stack[Bp], raise an error. Operand and stack[Bp] must be ints",
	"pushparm": "Push function parameter onto stack. Push the parameter index contained in the operand, numbered from 0. Parameter index 0 is at Bp-1, index 1 is at Bp-2, etc. Operand must be an int.",
	"copys":    "Copy a segment of the stack to the end of the stack. The portion at op.Offset from the end till op.Offset-op.Len from the end is copied to the end. Order is preserved.",
	"halt":     "Halt the VM",
	"clone":    "Clone the top of the stack. Pop the stack, call Clone on it, and push the result.",
	"load":     "Load data slot into stack. Reads the data slot at operand and pushes it. Operand must be an int.",
	"store":    "Store top of stack to data slot. Pop the stack and store the value into the data slot at operand. Operand must be an int.",
	"stores":   "Store stack element to data slot at address on top of stack. Pop the top of the stack as data slot index, and then pop the top of the stack again and store it to the data slot at that index. Top of stack must be an int.",
	"tload":    "Load value from table. Pop the table on the top of the stack, retrieve the element from it at the index given by the operand, and push that element. Top of stack must be a TableOperand, and operand must be an int.",
	"tstore":   "Store value to table. Pop the element on the top of the stack, then store it in the table on the new top of the stack at the index given by the operand. The table is not popped. ",
	"tmake":    "Make a table. The N topmost stack elements are popped and stored in a table, where N is the operand, and the table is then pushed onto the stack. The operand must be an int.",
	"alloc":    "Allocate a new data slot. The index of the new slot is pushed on the stack.",
	"free":     "Free a data slot previously allocated with alloc. The slot to be freed is popped from the top of the stack.",
	"mkclsr":   "Make a closure. The top of the stack is popped and stored in the closure as the closure environment. The lambda address of the closure is taken from the operand. The resulting closure is pushed on the stack.",
}

var InvalidOperandType = fmt.Errorf("operand is not valid for instruction")
var InvalidTopOfStackType = fmt.Errorf("value on top of stack is not valid for instruction")
var InvalidOperandValue = fmt.Errorf("operand value is out of range")
var InvalidFunction = fmt.Errorf("no function with that index")
var InvalidStackSize = fmt.Errorf("stack size is not valid for the instruction")
var InvalidAddress = fmt.Errorf("address into data segment is invalid")
var InvalidVariableType = fmt.Errorf("variable had wrong type")
var InvalidArgumentCount = fmt.Errorf("wrong number of arguments passed to function")
var InvalidGeneralRegister = fmt.Errorf("no such general register is defined")

// push an immediate value onto the stack
func pushOpHandler(state *vm.State, i *vm.Instruction) error {
	state.Stack.Push(i.Operand)
	return nil
}

// handle 'load' instruction.  Parameter to the instruction
// is a variable index (index into the data segment)
// pushes the value of the variable
func loadOpHandler(state *vm.State, i *vm.Instruction) error {
	ptr, ok := i.Operand.(Ref)
	if !ok {
		return InvalidOperandType
	}

	val, ok := state.GetData(int(ptr))
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
	ptr, ok := i.Operand.(Ref)
	if !ok {
		return InvalidOperandType
	}

	val := state.Stack.Pop()

	ok = state.SetData(int(ptr), val)
	if !ok {
		// restore stack
		state.Stack.Push(val)
		return InvalidAddress
	}

	return nil
}

// handle 'stores' instruction (store to adress from stack.
// top of stack is the address to store into, value is
// second
func storeStackOpHandler(state *vm.State, i *vm.Instruction) error {
	top := state.Stack.Pop()

	ptr, ok := top.(Ref)
	if !ok {
		state.Stack.Push(top)
		return InvalidTopOfStackType
	}

	val := state.Stack.Pop()

	ok = state.SetData(int(ptr), val)
	if !ok {
		// restore stack
		state.Stack.Push(val)
		state.Stack.Push(top)
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
	arg, ok := i.Operand.(Ref)
	if !ok {
		return InvalidOperandType
	}

	state.Stack.Push(state.Ip)
	state.Ip = int(arg - 1)
	return nil
}

// handle 'calli' instruction (indirect call). Parameter to the instruction
// is a variable index (index into the data segment)
// sets Ip (instruction pointer) to the address of the call - 1, and pushes return address
func callIndirectOpHandler(state *vm.State, i *vm.Instruction) error {
	ptr, ok := i.Operand.(Ref)
	if !ok {
		return InvalidOperandType
	}

	argIntf, ok := state.GetData(int(ptr))
	if !ok {
		return InvalidAddress
	}

	var newIp int
	switch t := argIntf.(type) {
	case int:
		newIp = t
	case Ref:
		newIp = int(t)
	case LambdaClosureOperand:
		newIp = int(t.LambdaAddr)
		tbl, ok := state.GetData(int(t.ClosureEnv))
		if !ok {
			return InvalidAddress
		}
		state.Gr[ClosureEnvGenReg] = tbl

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
	case Ref:
		newIp = int(t)
	case LambdaClosureOperand:
		newIp = int(t.LambdaAddr)
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
		return InvalidTopOfStackType
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
	if index >= len(state.Stack) || index < 0 {
		return InvalidOperandValue
	}
	state.Stack.Push(state.Stack[index])
	return nil
}

// swap top two elements of the stack
func swapOpHandler(state *vm.State, i *vm.Instruction) error {
	if len(state.Stack) < 2 {
		return InvalidStackSize
	}
	state.Stack.Swap()
	return nil
}

func haltOpHandler(state *vm.State, i *vm.Instruction) error {
	state.Halt = true
	return nil
}

func cloneOpHandler(state *vm.State, i *vm.Instruction) error {
	top := state.Stack.Pop()
	val, ok := Clone(top)
	if !ok {
		// might be an internal type
		val, ok = internalClone(top)
	}
	if !ok {
		state.Stack.Push(top)
		return InvalidTopOfStackType
	}
	state.Stack.Push(val)
	return nil
}

// Copy a segment of the stack to the end of the stack
// The portion at op.Offset from the end till op.Offset-op.Len from
// the end is copied to the end. Order is preserved.
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

// handle 'tload' instruction. Load the value from the table
// at top of stack, at index in operand. push it on the stack.
func tloadOpHandler(state *vm.State, i *vm.Instruction) error {
	//top := state.Stack.Top()
	top := state.Stack.Pop()
	tbl, ok := top.(Table)
	if !ok {
		state.Stack.Push(top)
		return InvalidTopOfStackType
	}

	ndx, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}

	val := tbl[ndx]
	if !ok {
		return InvalidAddress
	}

	state.Stack.Push(val)
	return nil
}

func tstoreOpHandler(state *vm.State, i *vm.Instruction) error {
	val := state.Stack.Pop()

	top := state.Stack.Top()
	tbl, ok := top.(*Table)
	if !ok {
		return InvalidTopOfStackType
	}

	ndx, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}

	(*tbl)[ndx] = val
	if !ok {
		return InvalidAddress
	}

	return nil
}

func tmakeOpHandler(state *vm.State, i *vm.Instruction) error {
	cnt, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}

	tbl := make(Table, cnt)

	for i := 0; i < cnt; i++ {
		val := state.Stack.Pop()
		tbl[i] = val
	}
	state.Stack.Push(tbl)
	return nil
}

func allocOpHandler(state *vm.State, i *vm.Instruction) error {
	state.Stack.Push(Ref(state.AllocData()))
	return nil
}

func freeOpHandler(state *vm.State, i *vm.Instruction) error {
	slot, ok := state.Stack.Pop().(Ref)
	if !ok {
		return InvalidOperandType
	}
	state.FreeData(int(slot))
	return nil
}

func makeClosureOpHandler(state *vm.State, i *vm.Instruction) error {
	var val LambdaClosureOperand

	// need three inputs:
	// - lambda address
	// - closure tbl address

	lambda, ok := i.Operand.(Ref)
	if !ok {
		return InvalidOperandType
	}

	top := state.Stack.Pop()

	ptr, ok := top.(Ref)
	if !ok {
		state.Stack.Push(top)
		return InvalidOperandType
	}

	val.ClosureEnv = ptr
	val.LambdaAddr = lambda

	state.Stack.Push(val)

	return nil
}

func setGenRegisterOpHandler(state *vm.State, i *vm.Instruction) error {
	top := state.Stack.Pop()
	ndx, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}
	if ndx < 0 || ndx >= len(state.Gr) {
		return InvalidGeneralRegister
	}
	state.Gr[ndx] = top
	return nil
}

func getGenRegisterOpHandler(state *vm.State, i *vm.Instruction) error {
	ndx, ok := i.Operand.(int)
	if !ok {
		return InvalidOperandType
	}
	if ndx < 0 || ndx >= len(state.Gr) {
		return InvalidGeneralRegister
	}

	state.Stack.Push(state.Gr[ndx])
	return nil
}
