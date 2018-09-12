package compiler

import "sort"

type envKey struct {
	fnName  string
	parmNdx int
}

type envVal struct {
	fnAndParamIndex
	id int
}

type env map[envKey]envVal

type closure struct {
	nextId int
	// These are the parameters of ancestor fns that the lambda refers
	// to. The fn creating the lambda needs to put these in a table.
	env env
}

func newClosure() *closure {
	return &closure{0, make(env)}
}

func (c *closure) addEnvEntry(f *fnAndParamIndex) (id int) {
	key := envKey{f.node.Name, f.parm}
	if v, ok := c.env[key]; ok {
		return v.id
	}
	v := envVal{*f, c.nextId}
	c.nextId++
	c.env[key] = v
	return v.id
}

// sortedEntries returns the closures env entries
// in descending order of id
func (c *closure) sortedEntries() (vals []envVal) {
	vals = make([]envVal, len(c.env))
	for _, v := range c.env {
		vals = append(vals, v)
	}
	// sort reversed
	sort.Slice(vals, func(i, j int) bool {
		return vals[i].id >= vals[j].id
	})
	return
}
