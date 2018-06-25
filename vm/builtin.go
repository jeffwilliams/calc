package vm

import "fmt"

var InvalidBuiltinError = fmt.Errorf("Invalid builtin")

type BuiltinSet interface {
	Func(index int) (Func, error)
	Name(index int) string
}

type BuiltinDescr struct {
	//Id      int
	Name string
	Func Func
}

type BuiltinTable struct {
	Funcs []Func
	Names map[int]string
}

func (t BuiltinTable) Func(index int) (Func, error) {
	if index >= len(t.Funcs) {
		return nil, InvalidBuiltinError
	}
	return t.Funcs[index], nil
}

func (t BuiltinTable) Name(index int) string {
	if index >= len(t.Funcs) {
		return "invalid"
	}
	return t.Names[index]
}

func NewBuiltinTable(descTable []BuiltinDescr) (bt *BuiltinTable, err error) {
	bt = &BuiltinTable{
		Funcs: make([]Func, len(descTable)),
		Names: map[int]string{},
	}

	for i, v := range descTable {
		bt.Funcs[i] = v.Func
		bt.Names[i] = v.Name
	}
	return
}
