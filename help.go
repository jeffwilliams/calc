package main

import (
	"fmt"
	"sort"
)

type NumParamer interface {
	NumParams() int
}

func printFuncHelp() {
	keys := make([]string, len(Funcs))
	i := 0
	for k := range Funcs {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := Funcs[k]
		p, ok := v.(NumParamer)
		if ok {
			fmt.Printf("%s(", k)
			num := p.NumParams()
			if num >= 0 {
				for i := 0; i < p.NumParams(); i++ {
					if i > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("p%d", i+1)
				}
			} else {
				fmt.Printf("...")
			}
			fmt.Printf("): %s\n", v.Help())
		} else {
			fmt.Printf("%s: %s\n", k, v.Help())
		}
	}
}
