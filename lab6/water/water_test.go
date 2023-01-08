package water

import (
	"fmt"
	"testing"
)

var waterTests = []struct {
	input            string
	wantNumMolecules int
}{
	{"OHH", 1},
	{"HHO", 1},
	{"HOH", 1},
	{"OOHHHH", 2},
	{"HHOHHO", 2},
	{"HOHHHO", 2},
	{"OHHHHO", 2},
	{"HHOHOH", 2},
	{"HOHHOH", 2},
	{"OHHHOH", 2},
	{"HHOOHH", 2},
	{"HOHOHH", 2},
	{"OHHOHH", 2},
	{"HHHHHHOOO", 3},
	{"OOOHHHHOHHHH", 4},
	{"HHHOOHOHHOHHHHO", 5},
	{"HOOHHHHOOOOHHHHOHHHHH", 7},
	{"HHHOHOHOHHOHHHHOOOHHHOOHOHHHHH", 10},
}

func TestWater(t *testing.T) {
	for _, test := range waterTests {
		water := New()
		gotMolecules := water.Make(test.input)
		gotNumMolecules := water.Molecules()
		if gotNumMolecules != test.wantNumMolecules {
			t.Errorf("water.Make(%s) = %d molecules, expected %d", test.input, gotNumMolecules, test.wantNumMolecules)
			continue
		}
		if !isWaterSequence(gotMolecules) {
			t.Errorf("water.Make(%s) = %q is invalid", test.input, gotMolecules)
		}
	}
}

func Example_one() {
	water := New()
	water.Make("HOH")
	fmt.Println(water.Molecules())
	// Output: 1
}

func Example_two() {
	water := New()
	water.Make("OOHHHH")
	fmt.Println(water.Molecules())
	// Output: 2
}
