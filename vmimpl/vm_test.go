package vmimpl

import (
	"fmt"
	"github.com/jeffwilliams/calc/vm"
	"math/big"
	"testing"
)

func TestProgs(t *testing.T) {

	m, err := NewVM()
	if err != nil {
		t.Fatalf("creating vm failed: %v", err)
	}

	tests := []struct {
		name     string
		program  []vm.Instruction
		expected *big.Int
	}{
		{
			"add_2",
			[]vm.Instruction{
				vm.Instruction{opcode(pushOp), big.NewInt(1)},
				vm.Instruction{opcode(pushOp), big.NewInt(2)},
				vm.Instruction{opcode(iAddOp), nil},
			},
			big.NewInt(3),
		},
		{
			"add_3",
			[]vm.Instruction{
				vm.Instruction{opcode(pushOp), big.NewInt(1)},
				vm.Instruction{opcode(pushOp), big.NewInt(2)},
				vm.Instruction{opcode(pushOp), big.NewInt(1)},
				vm.Instruction{opcode(iAddOp), nil},
				vm.Instruction{opcode(iAddOp), nil},
			},
			big.NewInt(4),
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			err = m.Run(tc.program)
			if err != nil {
				t.Fatalf("execution error: %v", err)
			}

			if len(m.State().Stack) == 0 {
				t.Fatalf("stack is empty after program")
			}

			if m.State().Stack.Top().(*big.Int).Cmp(tc.expected) != 0 {
				t.Fatalf("expected %v but got %v", tc.expected, m.State().Stack.Top())
			}

			fmt.Printf("%v\n", m.State().Stack.Top())
		})
	}
}
