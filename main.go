package main

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jeffwilliams/calc/ast"
	"github.com/jeffwilliams/calc/compiler"
	"github.com/jeffwilliams/calc/vm"
	"github.com/jeffwilliams/calc/vmimpl"
	flag "github.com/spf13/pflag"
)

//go:generate sh -c "$GOPATH/bin/pigeon calc.peg > gen_calc.go"
//go:generate $GOPATH/bin/genny -in eval.genny -out gen_eval.go gen "Number=big.Int,big.Float"
//go:generate $GOPATH/bin/genny -in op.genny -out gen_op.go gen "Op=add,sub,mul,quo,exp,and,or,lt,lte,gt,gte,eql"
//go:generate $GOPATH/bin/genny -in unary_op.genny -out gen_unary_op.go gen "Op=not,neg"

var outputBase numberBase = decimalBase

var completer = readline.NewPrefixCompleter()

func updateAutocomplete() {
	var items []readline.PrefixCompleterInterface

	for k := range GlobalVars {
		items = append(items, readline.PcItem(k))
	}

	for k := range Funcs {
		items = append(items, readline.PcItem(k))
	}

	var settings []readline.PrefixCompleterInterface
	for k := range Settings {
		settings = append(settings, readline.PcItem(k))
	}
	setItem := readline.PcItem("set", settings...)
	items = append(items, setItem)

	items = append(items, readline.PcItem("def"))
	items = append(items, readline.PcItem("help"))

	completer.SetChildren(items)
}

func LoadInitScript() (err error) {
	path := os.ExpandEnv("$HOME/.calcrc")

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	rdr := bufio.NewReader(file)
	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")

		_, err = Parse("init script", []byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	return
}

func printIntList(l BigIntList) {
	fmt.Printf("[")
	for i, v := range l {
		if i != 0 {
			fmt.Printf(", ")
		}
		fmt.Printf(outputBase.format(v))
	}
	fmt.Printf("]\n")
}

func printResult(parsed interface{}) {
	switch t := parsed.(type) {
	case *big.Int:
		fmt.Println(outputBase.format(parsed.(fmt.Formatter)))
	case *big.Float:
		fmt.Printf("%f\n", parsed)
	case []interface{}:
		for _, e := range t {
			printResult(e)
		}
	case BigIntList:
		printIntList(t)
	case BigFloatList:
		fmt.Printf("%s\n", parsed)
	case string:
		fmt.Printf("%s\n", parsed)
	default:
		// Don't print the results of statements
		//fmt.Println(parsed)
	}
}

var sharedCode *compiler.Shared

// link links the shared code with the passed shared object and returns the final code.
func linkWithShared(obj *compiler.Compiled) (code []vm.Instruction, err error) {
	if sharedCode != nil {
		fmt.Printf("linkWithShared: sharedCode before linking: %v\n", sharedCode)
		fmt.Printf("linkWithShared: obj.Shared before linking: %v\n", obj.Shared)
		obj.Shared.Link(sharedCode)
		fmt.Printf("linkWithShared: obj.Shared after linking: %v\n", obj.Shared)
	}
	sharedCode = &obj.Shared
	return obj.Linked()
}

func init() {
	vmimpl.Clone = clone
}

func main() {
	flag.VarP(&outputBase, "obase", "o", "Output number base. One of decimal, hex, integer. May be partial string.")
	flag.Parse()

	input := strings.Join(flag.Args(), " ")

	LoadInitScript()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		AutoComplete: completer,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: loading readline failed: %v\n", err)
		return
	}

	vmach, builtinIndexes, err := NewVM()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: creating VM failed: %v\n", err)
		return
	}

	var line string

	for {
		if input != "" {
			line = input
			input = ""
		} else {
			line, err = rl.Readline()
			if err != nil {
				if err == readline.ErrInterrupt {
					break
				}
				if err == io.EOF && !readline.IsTerminal(readline.GetStdin()) {
					break
				}
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
		}

		parsed, err := Parse("last line", []byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		} else {
			SetGlobal("last", parsed)
		}

		fmt.Printf("parsed type: %T\n", parsed)
		if t, ok := parsed.(*ast.Stmts); ok {
			t.GetMeta()
			f := func(node interface{}, depth int) bool {
				d := depth
				for ; depth > 0; depth-- {
					fmt.Printf("  ")
				}
				fmt.Printf("Node: %T %+v (depth %d)\n", node, node, d)
				return true
			}
			ast.Walk(f, ast.Pre, t)

			obj, err := compiler.Compile(t, builtinIndexes, sharedCode)
			if err != nil {
				fmt.Printf("compilation failed: %v\n", err)
				continue
			}

			code, err := linkWithShared(obj)
			if err != nil {
				fmt.Printf("final link failed: %v\n", err)
				continue
			}

			fmt.Printf("compiled code:\n")
			for i, instr := range code {
				fmt.Printf("%d: %s\n", i, vmach.InstructionString(&instr))
			}

			fmt.Printf("shared code:\n")
			fmt.Printf("%+v\n", sharedCode)

			stepFunc := func(state *vm.State) {
				fmt.Println(state.Summary(vmach.InstructionSet()))
			}

			vmopts := vm.RunOpts{StepFunc: stepFunc}

			err = vmach.Run(code, &vmopts)
			if err != nil {
				fmt.Printf("execution failed: %v\n", err)
				if e, ok := err.(vm.ExecError); ok {
					fmt.Printf("%s\n", e.Details())
				}
				continue
			}

			if len(vmach.State().Stack) == 0 {
				printResult("Error: stack is empty after execution")
				continue
			}

			fmt.Printf("%v\n", vmach.State().Stack.Top())
		}

		printResult(parsed)
	}
}
