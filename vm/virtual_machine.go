package vm

import "fmt"

type State struct {
	Stack Stack
	halt  bool
	ip    int
}

type OpcodeHandler func(state *State, i *Instruction) error

type VM struct {
	is    InstructionSet
	state State
}

func NewVM(is InstructionSet) (vm *VM, err error) {
	vm = &VM{
		state: State{
			Stack: make(Stack, 0, 10000),
		},
		is: is,
	}
	return
}

func (vm *VM) Run(prog []Instruction) error {
	vm.state.ip = 0
	vm.state.Stack = []interface{}{}
	for !vm.state.halt && vm.state.ip < len(prog) {
		i := prog[vm.state.ip]
		//fmt.Printf("Calling opcode %s at ip=%d\n", vm.is.Name(i.Opcode), vm.state.ip)

		h, err := vm.is.Handler(i.Opcode)
		if err != nil {
			return fmt.Errorf("error at ip=%d: %v", vm.state.ip, err)
		}

		h(&vm.state, &i)
		if err != nil {
			return fmt.Errorf("error at ip=%d: %v", vm.state.ip, err)
		}
		vm.state.ip++
	}
	return nil
}

func (vm *VM) State() *State {
	return &vm.state
}
