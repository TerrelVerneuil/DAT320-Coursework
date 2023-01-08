//go:build !race

package stack

import (
	"testing"
)

var stackTests = []struct {
	name     string
	ops      []string
	args     []string
	wantSize int
	wantElem interface{}
}{
	{"Initial Empty Stack", []string{"Size"}, []string{""}, 0, ""},
	{"Push One Element", []string{"Push"}, []string{"Item1"}, 1, ""},
	{"Pop One Element", []string{"Pop"}, []string{""}, 0, "Item1"},
	{"Push Three Items", []string{"Push", "Push", "Push"}, []string{"i1", "k2", "x4"}, 3, ""},
	{"Pop One Item", []string{"Pop"}, []string{""}, 2, "x4"},
	{"Pop Another Item", []string{"Pop"}, []string{""}, 1, "k2"},
	{"Pop One More Item", []string{"Pop"}, []string{""}, 0, "i1"},
	{"Push Two New Element", []string{"Push", "Push"}, []string{"x", "y"}, 2, ""},
	{"Pop New Element", []string{"Pop"}, []string{""}, 1, "y"},
	{"Pop The Other New Element", []string{"Pop"}, []string{""}, 0, "x"},
	{"Pop From Empty", []string{"Pop"}, []string{""}, 0, nil},
	{"Pop From Empty Again", []string{"Pop"}, []string{""}, 0, nil},
}

func TestStackOps(t *testing.T) {
	stacks := map[string]Stack{
		"UnsafeStack": new(UnsafeStack),
		"SafeStack":   new(SafeStack),
		"SliceStack":  NewSliceStack(),
		"CspStack":    NewCspStack(),
	}
	for name, stack := range stacks {
		t.Run(name, func(t *testing.T) {
			for _, test := range stackTests {
				for index, action := range test.ops {
					switch action {
					case "Size":
						if got := stack.Size(); got != test.wantSize {
							t.Errorf("%25s: %s(%s) = %v, want: %v", test.name, action, test.args[index], got, test.wantSize)
						}
					case "Push":
						stack.Push(test.args[index])
						if got := stack.Pop(); got != test.args[index] {
							t.Errorf("%25s: %s(%s); Pop() = %v, want: %v", test.name, action, test.args[index], got, test.args[index])
						}
						stack.Push(test.args[index])
					case "Pop":
						if got := stack.Pop(); got != test.wantElem {
							t.Errorf("%25s: %s(%s) = %v, want: %v", test.name, action, test.args[index], got, test.wantElem)
						}
					}
				}
				if got := stack.Size(); got != test.wantSize {
					t.Errorf("%25s: Size() = %v, want: %v", test.name, got, test.wantSize)
				}
			}
		})
	}
}

func BenchmarkStacks(b *testing.B) {
	const numOps = 10000
	stacks := map[string]Stack{
		"SafeStack":  new(SafeStack),
		"SliceStack": NewSliceStack(),
		"CspStack":   NewCspStack(),
	}
	for name, stack := range stacks {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for i := 0; i < numOps; i++ {
					stack.Push(i)
				}
				for j := 0; j < numOps; j++ {
					stack.Pop()
				}
			}
		})
	}
}
