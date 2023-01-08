//go:build race

package stack

import (
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	Len = iota
	Push
	Pop
)

func TestConcurrentStackAccess(t *testing.T) {
	const numOps = 10_000

	// seed the random number generator with the time since the Unix epoch
	rand.Seed(time.Now().Unix())
	numGoroutines := runtime.NumCPU()
	t.Logf("Running with %d goroutines (= number of CPUs)", numGoroutines)

	stacks := map[string]Stack{
		"SafeStack":  new(SafeStack),
		"SliceStack": NewSliceStack(),
		"CspStack":   NewCspStack(),
	}
	for name, stack := range stacks {
		t.Run(name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(numGoroutines)
			for i := 0; i < numGoroutines; i++ {
				go func(i int) {
					var cnt int
					for j := 0; j < numOps; j++ {
						op := rand.Intn(3)
						switch op {
						case Len:
							stack.Size()
						case Push:
							stack.Push("Data" + strconv.Itoa(i) + strconv.Itoa(j))
						case Pop:
							_ = stack.Pop()
						}
						cnt = j + 1
					}
					if stack.Size() == -1 {
						t.Errorf("CspStack.Size() = %d, expected >= 0", stack.Size())
					} else {
						t.Logf("G%02d: ops: %d, remaining items on stack: %d\n", i, cnt, stack.Size())
					}
					defer wg.Done()
				}(i)
			}
			wg.Wait()
		})
	}
}
