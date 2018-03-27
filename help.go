package main

import "fmt"

type NumParamer interface {
	NumParams() int
}

func printFuncHelp() {
	for k, v := range Funcs {
		p, ok := v.(NumParamer)
		if ok {
			fmt.Printf("%s(", k)
			for i := 0; i < p.NumParams(); i++ {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("p%d", i+1)
			}
			fmt.Printf("): %s\n", v.Help())
		} else {
			fmt.Printf("%s: %s\n", k, v.Help())
		}
	}
}
