package main

import (
	"fmt"
	"math/big"

	"github.com/jeffwilliams/calc/vm"
	. "github.com/jeffwilliams/calc/vmimpl"
)

func main() {
	/*
		m, err := NewVM()
		if err != nil {
			fmt.Printf("creating vm failed: %v", err)
			return
		}

		prog := []vm.Instruction{
			// Call "add 1" with argument 5
			I(PushOp, big.NewInt(5)),
			I(CallOp, 5),

			// Continue after call and add 2 to result of call
			I(PushOp, big.NewInt(2)),
			I(IAddOp, nil),
			I(HaltOp, nil),

			// function "add 1" that adds 1 to the passed parameter
			I(EnterOp, nil),
			I(PushOp, big.NewInt(1)),
			I(PushParmOp, 0),
			I(IAddOp, nil),
			I(LeaveOp, 1),
			I(ReturnOp, nil),
		}

		runOpts := vm.RunOpts{StepFunc: func(state *vm.State) {
			fmt.Printf("%s\n", state.Summary(m.InstructionSet()))

		}}

		err = m.Run(prog, &runOpts)
		if err != nil {
			fmt.Printf("execution error: %v\nexecution state at error:\n%v", err, err.(vm.ExecError).Details())
		}
	*/
}
