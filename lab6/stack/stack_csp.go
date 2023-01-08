package stack

type stackOperation int

const (
	length stackOperation = iota
	push
	pop
)

type stackCommand struct {
	op     stackOperation
	// TODO(student) Add necessary fields
}

// CspStack is a struct with methods needed to implement the Stack interface.
type CspStack struct {
	size int
	// TODO(student) Add necessary fields
}

// NewCspStack returns an empty CspStack.
func NewCspStack() *CspStack {
	// TODO(student) Implement constructor and start handling commands
	return &CspStack{}
}

// Size returns the size of the stack.
func (cs *CspStack) Size() int {
	// TODO(student) Implement size
	return -1
}

// Push pushes value onto the stack.
func (cs *CspStack) Push(value interface{}) {
	// TODO(student) Implement push
}

// Pop pops the value at the top of the stack and returns it.
func (cs *CspStack) Pop() (value interface{}) {
	// TODO(student) Implement pop
	return nil
}

func (cs *CspStack) run() {
	// TODO(student) Implement handlers for each stack command
}
