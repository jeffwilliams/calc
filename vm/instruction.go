package vm

import "fmt"

type Instruction struct {
	// opcode is an index into the opcode table
	Opcode uint8
	// operands.
	Operand interface{}

	// index to an immediate big.Int value
	// index into a variable table
	// index into a function table
	// ??
}

var InvalidOpcodeError = fmt.Errorf("Invalid opcode")

type InstructionSet interface {
	Handler(opcode uint8) (OpcodeHandler, error)
	Name(opcode uint8) string
}

type InstructionTable struct {
	// Handlers maps an instruction opcode to its handler
	Handlers []OpcodeHandler
	// Names maps an instruction opcode to its name
	Names map[uint8]string
	// Opcode maps an instruction description's id to its assigned opcode
	//Opcode map[int]uint8
}

func (i InstructionTable) Handler(opcode uint8) (OpcodeHandler, error) {
	if int(opcode) >= len(i.Handlers) {
		return nil, InvalidOpcodeError
	}
	return i.Handlers[opcode], nil
}

func (i InstructionTable) Name(opcode uint8) string {
	if int(opcode) >= len(i.Handlers) {
		return "invalid"
	}
	return i.Names[opcode]
}

type InstructionDescr struct {
	//Id      int
	Name    string
	Handler OpcodeHandler
}

func NewInstructionTable(descTable []InstructionDescr) (it *InstructionTable, err error) {
	it = &InstructionTable{
		Handlers: make([]OpcodeHandler, len(descTable)),
		Names:    map[uint8]string{},
		//Opcode:   map[int]uint8{},
	}

	if len(descTable) > 256 {
		err = fmt.Errorf("At most 255 opcodes are supported")
		return
	}

	for i, v := range descTable {
		it.Handlers[i] = v.Handler
		it.Names[uint8(i)] = v.Name
		//it.Opcode[v.Id] = uint8(i)
	}
	return
}
