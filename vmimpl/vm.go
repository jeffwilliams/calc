package vmimpl

import (
	"fmt"

	"github.com/jeffwilliams/calc/vm"
)

//var opcodes []OpcodeHandler
var InstructTable *vm.InstructionTable

var instructionDesc = []vm.InstructionDescr{
	{"invalid", haltOpHandler},
	{"iadd", iAddOpHandler},
	{"push", pushOpHandler},
	{"pop", popOpHandler},
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
	{"tload", tloadOpHandler},
	{"tstore", tstoreOpHandler},
	{"tmake", tmakeOpHandler},
}

var instructionOpcode = map[string]uint8{}

func opcode(name string) (code uint8, err error) {
	var ok bool
	code, ok = instructionOpcode[name]
	if !ok {
		err = fmt.Errorf("no such instruction '%s'", name)
		return
	}
	return
}

// Convenience function for building an instruction from an instruction id and operand
func I(name string, operand interface{}) vm.Instruction {
	code, err := opcode(name)
	if err != nil {
		panic(err)
	}
	return vm.Instruction{code, operand}
}

func init() {
	var err error
	InstructTable, err = vm.NewInstructionTable(instructionDesc)
	if err != nil {
		panic(err)
	}

	for i, v := range instructionDesc {
		instructionOpcode[v.Name] = uint8(i)
	}
}

func NewVM(bs vm.BuiltinSet) (m *vm.VM, err error) {
	if Clone == nil {
		err = fmt.Errorf("Error: the vmimpl.Clone function pointer is not set. Please set it to main.clone.")
		return
	}
	m, err = vm.NewVM(InstructTable, bs)
	return
}
