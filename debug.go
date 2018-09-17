package main

import (
	"strings"
)

type DbgFlags int

const (
	DbgFlagParse DbgFlags = 1 << iota
	DbgFlagAst
	DbgFlagVm
	DbgFlagAsmHelp
	DbgFlagMax
)

func parseDebugFlags(s string) DbgFlags {
	if strings.ContainsRune(s, 'a') {
		return DbgFlagMax - 1
	}

	var flags DbgFlags

	if strings.ContainsRune(s, 'p') {
		flags |= DbgFlagParse
	}
	if strings.ContainsRune(s, 's') {
		flags |= DbgFlagAst
	}
	if strings.ContainsRune(s, 'v') {
		flags |= DbgFlagVm
	}
	if strings.ContainsRune(s, 'h') {
		flags |= DbgFlagAsmHelp
	}
	return flags
}
