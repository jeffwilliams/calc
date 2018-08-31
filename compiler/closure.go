package compiler

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
