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

// CopyToEnd copies a segment of the stack starting at `depth` from the end of length `len`
// to the end of the stack
func (s *Stack) CopyToEnd(depth, length int) {

	if depth+length > len(*s) {
		length = len(*s) - depth
		if length <= 0 {
			return
		}
	}

	*s = append(*s, (*s)[len(*s)-depth-length:len(*s)-depth]...)
}
