package stack

// TODO(student) Add necessary fields and synchronization primitives

// SafeStack holds the top element of the stack and its size.
type SafeStack struct {
	top  *Element
	size int
}

// Size returns the size of the stack.
func (ss *SafeStack) Size() int {
	return ss.size
}

// Push pushes value onto the stack.
func (ss *SafeStack) Push(value interface{}) {
	ss.top = &Element{value, ss.top}
	ss.size++
}

// Pop pops the value at the top of the stack and returns it.
func (ss *SafeStack) Pop() (value interface{}) {
	if ss.size > 0 {
		value, ss.top = ss.top.value, ss.top.next
		ss.size--
		return
	}
	return nil
}
