package fizzbizz

import (
	"flag"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var fizzBizzTests = []struct {
	n      int
	result string
}{
	{1, "1"},
	{2, "12"},
	{3, "12Fizz"},
	{4, "12Fizz4"},
	{5, "12Fizz4Bizz"},
	{6, "12Fizz4BizzFizz"},
	{7, "12Fizz4BizzFizz7"},
	{8, "12Fizz4BizzFizz78"},
	{9, "12Fizz4BizzFizz78Fizz"},
	{10, "12Fizz4BizzFizz78FizzBizz"},
}

var userInput = flag.Int("max", 20, "max value")

func TestFizzBizz(t *testing.T) {
	for _, test := range fizzBizzTests {
		gotResult := FizzBizz(test.n)
		if diff := cmp.Diff(test.result, gotResult); diff != "" {
			t.Errorf("fizzbizz.FizzBizz(%d): (-want +got):\n%s", test.n, diff)
		}
	}
}

func TestFizzBizzWithUserInput(t *testing.T) {
	fmt.Println(FizzBizz(*userInput))
}

func Example_ten() {
	fmt.Println(FizzBizz(10))
	// Output: 12Fizz4BizzFizz78FizzBizz
}
