package vm

type Stack []interface{}

func (s *Stack) Top() interface{} {
	return (*s)[len(*s)-1]
}

func (s *Stack) Pop() (r interface{}) {
	if len(*s) > 0 {
		r = (*s)[len(*s)-1]
		*s = (*s)[0 : len(*s)-1]
	}
	return
}

func (s *Stack) Push(v interface{}) {
	*s = append(*s, v)
}

func (s *Stack) Swap() {
	l := len(*s)
	if l >= 2 {
		(*s)[l-2], (*s)[l-1] = (*s)[l-1], (*s)[l-2]
	}
}
