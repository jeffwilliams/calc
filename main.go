package main

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/chzyer/readline"
	flag "github.com/spf13/pflag"
)

//go:generate sh -c "$GOPATH/bin/pigeon calc.peg > gen_calc.go"
//go:generate $GOPATH/bin/genny -in eval.genny -out gen_eval.go gen "Number=big.Int,big.Float"
//go:generate $GOPATH/bin/genny -in op.genny -out gen_op.go gen "Op=add,sub,mul,quo,exp,and,or"
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

func main() {
	flag.VarP(&outputBase, "obase", "o", "Output number base. One of decimal, hex, integer. May be partial string.")
	flag.Parse()

	LoadInitScript()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		AutoComplete: completer,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: loading readline failed: %v\n", err)
		return
	}

	for {
		line, err := rl.Readline()
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

		parsed, err := Parse("last line", []byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		} else {
			SetGlobal("last", parsed)
		}

		printResult(parsed)
	}
}
