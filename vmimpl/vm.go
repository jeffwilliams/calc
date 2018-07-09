package vmimpl

import (
	"fmt"

	"github.com/jeffwilliams/calc/vm"
)

//var opcodes []OpcodeHandler
var instructTable *vm.InstructionTable

var instructionDesc = []vm.InstructionDescr{
	{"invalid", haltOpHandler},
	{"iadd", iAddOpHandler},
	{"push", pushOpHandler},
	{"pushi", pushIndirectOpHandler},
	{"callb", callBuiltinOpHandler},
	{"call", callOpHandler},
	{"calli", callIndirectOpHandler},
	{"pop", popOpHandler},
	{"return", returnOpHandler},
	{"enter", enterOpHandler},
	{"pushparm", pushParmOpHandler},
	{"leave", leaveOpHandler},
	{"halt", haltOpHandler},
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
	instructTable, err = vm.NewInstructionTable(instructionDesc)
	if err != nil {
		panic(err)
	}

	for i, v := range instructionDesc {
		instructionOpcode[v.Name] = uint8(i)
	}
}

func NewVM(bs vm.BuiltinSet) (m *vm.VM, err error) {
	m, err = vm.NewVM(instructTable, bs)
	return
}
