package vmimpl

import (
	"github.com/jeffwilliams/calc/vm"
)

//var opcodes []OpcodeHandler
var instructTable *vm.InstructionTable

var instructionDesc = []vm.InstructionDescr{
	{IAddOp, "iadd", iAddOpHandler},
	{PushOp, "push", pushOpHandler},
	{CallBOp, "callb", callBuiltinOpHandler},
	{CallOp, "call", callOpHandler},
	{PopOp, "pop", popOpHandler},
	{ReturnOp, "return", returnOpHandler},
	{EnterOp, "enter", enterOpHandler},
	{PushParmOp, "pushparm", pushParmOpHandler},
	//{SwapOp, "swap", swapOpHandler},
	{LeaveOp, "leave", leaveOpHandler},
	{HaltOp, "halt", haltOpHandler},
}

func opcode(instructionId int) uint8 {
	return instructTable.Opcode[instructionId]
}

// Convenience function for building an instruction from an instruction id and operand
func I(instructionId int, operand interface{}) vm.Instruction {
	return vm.Instruction{opcode(instructionId), operand}
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
