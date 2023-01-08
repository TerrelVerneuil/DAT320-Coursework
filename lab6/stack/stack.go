package stack

// UnsafeStack is a struct with methods needed to implement the Stack interface.
type UnsafeStack struct {
	top  *Element
	size int
}

// Size returns the size of the stack.
func (us *UnsafeStack) Size() int {
	return us.size
}

// Push pushes value onto the stack.
func (us *UnsafeStack) Push(value interface{}) {
	us.top = &Element{value, us.top}
	us.size++
}

// Pop pops the value at the top of the stack and returns it.
func (us *UnsafeStack) Pop() (value interface{}) {
	if us.size > 0 {
		value, us.top = us.top.value, us.top.next
		us.size--
		return
	}
	return nil
}
