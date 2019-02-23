package main

import (
	"fmt"

	"github.com/jeffwilliams/calc/vm"
	"github.com/jeffwilliams/calc/vmimpl"
)

// NewVM creates and returns a new VM. It constructs the set of builtin functions
// based on the Funcs array.
func NewVM() (m *vm.VM, builtinIndexes map[string]int, err error) {
	builtinDesc := []vm.BuiltinDescr{}

	for name, fn := range Funcs {
		builtinDesc = append(builtinDesc, vm.BuiltinDescr{name, fn})
	}

	builtinIndexes = map[string]int{}
	for i, v := range builtinDesc {
		builtinIndexes[v.Name] = i
	}

	builtinTable, err := vm.NewBuiltinTable(builtinDesc)
	if err != nil {
		panic(fmt.Sprintf("creating builtin table failed: %v", err))
	}

	m, err = vmimpl.NewVM(builtinTable)
	return
}

//func Compile(tree interface{}, builtinIndexes map[string]int) (code []vm.Instruction, err error) {
