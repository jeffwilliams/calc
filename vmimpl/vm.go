package vmimpl

import (
	"github.com/jeffwilliams/calc/vm"
)

//var opcodes []OpcodeHandler
var instructTable *vm.InstructionTable

var instructionDesc = []vm.InstructionDescr{
	{iAddOp, "iadd", iAddOpHandler},
	{pushOp, "push", pushOpHandler},
}

func opcode(instructionId int) uint8 {
	return instructTable.Opcode[instructionId]
}

func init() {
	var err error
	instructTable, err = vm.NewInstructionTable(instructionDesc)
	if err != nil {
		panic(err)
	}
}

func NewVM() (m *vm.VM, err error) {
	m, err = vm.NewVM(instructTable)
	return
}
