package core

type stack []uint16

func (s *stack) push(val uint16) {
	*s = append(*s, val)
}

func (s *stack) pop() uint16 {
	if len(*s) == 0 {
		return 0
	}
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}
