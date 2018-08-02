package vm

import (
	"bytes"
	"fmt"
)

type ExecError struct {
	Internal error
	VmState  State
	is       InstructionSet
}

func (e ExecError) Error() string {
	return e.Internal.Error()
}

func (e ExecError) Details() string {
	return e.VmState.Summary(e.is, nil, nil)
}

// StringWithStater defines objects that can return a String representing them
// if they are passed State.
type StringWithStater interface {
	StringWithState(s *State) string
}

type Func interface {
	Call(parms []interface{}) (result interface{}, err error)
}

type State struct {
	Stack    Stack
	Data     []interface{}
	Builtins BuiltinSet
	Prog     []Instruction
	//FuncParms  []interface{}
	Halt bool
	Ip   int
	// Base pointer
	Bp int
}

func (s State) InstructionString(is InstructionSet, instr *Instruction) string {
	var opString string
	if t, ok := instr.Operand.(StringWithStater); ok {
		opString = t.StringWithState(&s)
	} else {
		opString = fmt.Sprintf("%v", instr.Operand)
	}

	return fmt.Sprintf("%s %s (%T)", is.Name(instr.Opcode), opString, instr.Operand)
}

func (s State) Summary(is InstructionSet, codeNames, dataNames map[int]string) string {
	var buf bytes.Buffer

	printName := func(tbl map[int]string, i int) {
		if tbl != nil {
			if name, ok := tbl[i]; ok {
				fmt.Fprintf(&buf, "%s:\n", name)
			}
		}
	}

	fmt.Fprintf(&buf, "Ip: %d   Bp: %d\n", s.Ip, s.Bp)
	fmt.Fprintf(&buf, "Data:\n")
	for i := 0; i < len(s.Data); i++ {
		printName(dataNames, i)
		fmt.Fprintf(&buf, "  %d: %v (%T)\n", i, s.Data[i], s.Data[i])
	}
	fmt.Fprintf(&buf, "Stack:\n")
	for i := len(s.Stack) - 1; i >= 0; i-- {
		fmt.Fprintf(&buf, "  %d: %v (%T)\n", i, s.Stack[i], s.Stack[i])
	}
	if len(s.Stack) == 0 {
		fmt.Fprintf(&buf, "  (empty)\n")
	}
	fmt.Fprintf(&buf, "Instructions at Ip ± 10:\n")
	for i := s.Ip - 10; i < s.Ip+10 && i < len(s.Prog); i++ {
		if i < 0 {
			continue
		}

		printName(codeNames, i)

		instr := s.Prog[i]
		m := "  "
		if i == s.Ip {
			m = "=>"
		}

		fmt.Fprintf(&buf, "%s%d: %s\n", m, i, s.InstructionString(is, &instr))
	}
	return buf.String()
}

func (s State) GetData(i int) (val interface{}, ok bool) {
	if i < 0 || i >= len(s.Data) {
		return
	}

	val = s.Data[i]
	ok = true
	return
}

func (s *State) SetData(i int, val interface{}) (ok bool) {
	if i < 0 {
		return
	}

	if i >= len(s.Data) {
		nw := make([]interface{}, i+1)
		copy(nw, s.Data)
		s.Data = nw
	}

	s.Data[i] = val
	ok = true
	return
}

type OpcodeHandler func(state *State, i *Instruction) error

type VM struct {
	is    InstructionSet
	state State
}

func (vm VM) InstructionString(i *Instruction) string {
	return vm.state.InstructionString(vm.is, i)

}

func NewVM(is InstructionSet, bs BuiltinSet) (vm *VM, err error) {
	vm = &VM{
		state: State{
			Stack:    make(Stack, 0, 10000),
			Builtins: bs,
		},
		is: is,
	}
	return
}

type StepFunc func(state *State)

type RunOpts struct {
	StepFunc StepFunc
}

// If instructions generate an error, they must restore the stack (and other state) to the state previous
// to the instruction being executed, to help debugging.
func (vm *VM) Run(prog []Instruction, opts *RunOpts) error {
	vm.state.Ip = 0
	vm.state.Bp = 0
	vm.state.Stack = []interface{}{}
	vm.state.Prog = prog
	vm.state.Halt = false
	for !vm.state.Halt && vm.state.Ip < len(prog) {
		if opts != nil && opts.StepFunc != nil {
			opts.StepFunc(&vm.state)
		}
		ip := vm.state.Ip
		i := prog[vm.state.Ip]
		//fmt.Printf("Calling opcode %s at Ip=%d\n", vm.is.Name(i.Opcode), vm.state.Ip)

		h, err := vm.is.Handler(i.Opcode)
		if err != nil {
			//return fmt.Errorf("error at Ip=%d: %v", vm.state.Ip, err)
			return ExecError{err, vm.state, vm.is}
		}

		err = h(&vm.state, &i)
		if err != nil {
			if vm.state.Ip != ip {
				// restore Ip incase the instruction modified it.
				vm.state.Ip = ip
			}
			//return fmt.Errorf("error at Ip=%d: %v", vm.state.Ip, err)
			return ExecError{err, vm.state, vm.is}
		}
		vm.state.Ip++
	}
	return nil
}

func (vm *VM) State() *State {
	return &vm.state
}

func (vm *VM) InstructionSet() InstructionSet {
	return vm.is
}
