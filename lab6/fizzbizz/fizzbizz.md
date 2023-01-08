# Task: FizzBizz

This assignment aims to test your understanding of the `sync` package.
In this assignment, a number `max` is provided as the input.

You have to complete four goroutines and each one should process the numbers from in the range `[1,max]`.

The following goroutines must be completed according to these rules:

- `fizz()` that prints the word `Fizz`, if the number is divisible by 3 and not 5.
- `bizz()` that prints the word `Bizz`, if the number is divisible by 5 and not 3.
- `number()` that prints the `number` if it is not divisible by 3 and 5.
- `fizzBizz()` that prints the word `FizzBizz`, if the number is divisible by both 3 and 5.

The provided code contains additional instructions, including TODO items for you to complete.

You can check your implementation by running the test case:

```console
% go test -v -run TestFizzBizz
```

## Examples

With the inputs `max = 5` and `max = 15`, the output should be:

```console
% go test -v -run TestFizzBizzWithUserInput -max 5
12Fizz4Bizz
% go test -v -run TestFizzBizzWithUserInput -max 15
12Fizz4BizzFizz78FizzBizz11Fizz1314FizzBizz
```
