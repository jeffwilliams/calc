package compiler

type closure struct {
	baseIndex int
	// allocedNames is the parameters/variables already added to the closure and allocated in the var symbols.
	// The key is the actual name from the code, and the value is the internal generated name that was added
	// as the var symbol.
	// Values should be generated from compiler.closureNameGen
	allocedNames map[string]string
}

func newClosure() *closure {
	return &closure{allocedNames: make(map[string]string)}
}
