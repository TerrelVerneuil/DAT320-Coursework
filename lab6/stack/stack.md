# High-level Synchronization & Data Race Detection

In this part of the lab we will focus on high-level synchronization techniques using the Go programming language.

The Go language provides a built-in race detector that we will use to identify data races and verify implementations.
A *data race* occurs when two or more threads (or goroutines) access a variable concurrently and at least one of these accesses is a write.

We will work on a stack data structure that will be accessed concurrently from several goroutines.
The stack stores values of type `interface{}`, meaning any value type.
The stack interface is shown in Listing 1 and is found in the file `stack_iface.go`.
The interface contains three methods:

* `Size()` returns the current number of items on the stack,
* `Pop() interface{}` pops an item of the stack (`nil` if empty), while
* `Push(value interface{})` pushes an item onto the stack.

Listing 1: Stack interface

```go
type Stack interface {
    Size() int
    Push(value interface{})
    Pop() interface{}
}
```

For this lab we will use the tests defined in `stack_test.go` to verify the different stack implementations.
The tests can be run using the `go test` command. We will run one test at a time.
Running only a specific test can be achieved by supplying the `-run` flag together with a _regular expression_ indicating the test names.
For example, to only run the stack operations test for the `UnsafeStack` stack, use the command:

```console
go test -v -run TestStackOps/UnsafeStack
```

The `-v` flag enables verbose output.

There are two types of tests defined for each stack implementation we will be working on.
One test verifies a stack's operations, while the other is meant to test concurrent access using the race detector.
Study the test file for details.

As stated in the introduction, Go includes a built-in data race detector.
Read [Data Race Detector](<http://golang.org/doc/articles/race_detector.html>) for an introduction and usage examples.

## Task: Implement Thread-Safe Stacks

1. The file `stack_sync.go` is a copy of `stack.go`, but the type is renamed to `SafeStack`.
   Modify this file so that access to the `SafeStack` type is synchronized, i.e., can be safely accessed from concurrently running goroutines.

2. Check your implementation by running the `SafeStack` test with the data race detector enabled.
   The test should not produce any data race warnings.

   ```console
   go test -race -run TestConcurrentStackAccess/SafeStack
   ```

3. Go has a built-in high-level API for concurrent programming based on Communicating Sequential Processes (CSP).
   This API promotes synchronization through sending and receiving data via thread-safe channels, as opposed to traditional locking.
   The file `stack_csp.go` contains a `CspStack` type that implements the stack interface in Listing 1, but the actual method implementations are empty.
   The type also has a constructor function needed for this task.
   Modify this file so that access to the `CspStack` type is synchronized using Go's channels and goroutines.

   There is in this case an amount of overhead when using channels to achieve synchronization compared to locking.
   The main point for this task is to give an introduction on how to use channels (CSP) for synchronization.
   This will require some self-study if you are not familiar with Go’s CSP-based concurrent programming capabilities.
   A place to start can be the introduction found [here](http://golang.org/doc/effective_go.html#concurrency).

   To implement the `CspStack` using channels and goroutines, you must define a channel on which commands should be sent.
   The actual handling of commands must be done in the `run` method, which should be executed in a separate goroutine.
   The `run` method is expected to loop over the channel and process the commands in the order they are received.
   The `command` sent over the channel from the various `Stack` methods should contain the `type` of the command, the `value` to be stored in the stack, and the channel on which the `response` should be sent.

   Note that you should also ensure that the stack operations are implemented correctly.
   You can verify them by running:

   ```console
   go test -v -run TestStackOps/CspStack
   ```

4. Check your `CspStack` implementation by running the following test with the data race detector enabled.
   The test should not produce any data race warnings.

   ```console
   go test -v -race -run TestConcurrentStackAccess/CspStack
   ```

## CPU and Memory Profiling

In this part of the lab we will use a technique called profiling to dynamically analyze a program.
Profiling can among other things be used to measure an application’s CPU utilization and memory usage.
Being able to profile applications is very helpful when optimizing the performance of a program.
Profiling is an important part of *systems programming*.
This lab will give a short introduction to how profiling data can be analyzed.

Profiling for Go can be enabled through the `runtime/pprof` package or by using the testing package's profiling support.
Profiles can be analyzed and visualized using the `go tool pprof` program.

We will continue to use the stacks from the previous part of the lab.
The file `stack_test.go` contains the `BenchmarkStacks` function that can benchmark the three stack implementation.
The benchmark consists of pushing 10,000 items onto the stack and then popping them.
The stack implementations are not accessed concurrently so that the benchmarks are deterministic.
To run the benchmarks for the three stack implementations:

```console
go test -run none -bench BenchmarkStacks
```

Note that we provide `-run none` in this command, which doesn't match any tests in the `_test.go` file.

Read [Profiling Go Programs](https://blog.golang.org/pprof).
This blog post present a good introduction to Go's profiling abilities.
You should also examine the [testing](http://golang.org/pkg/testing/) package and [testing flags](http://golang.org/cmd/go/#Description_of_testing_flags) for information on how to run the benchmarks, and details about how Go's testing tool easily enables profiling when benchmarking.
Furthermore, we recommend reading [The Go Memory Model](http://golang.org/ref/mem) and [Introducing the Go Race Detector](http://blog.golang.org/race-detector).

### Task: Multiple Choice Questions About CPU and Memory Profiling and Data Races

Answer these multiple choice questions about [CPU and Memory Profiling and Data Races](profile_race_questions.md).

## Task: Benchmarking and Profiling Go Programs

In this task, you should fill in answers in the provided template: [`benchmark_report.md`](benchmark_report.md).
You can add figures in a directory `fig`, and add markdown links from the benchmark report file, so that the figures display nicely on GitHub's web page.

1. The file `stack_slice.go` contains a stack implementation, `SliceStack`, backed by a slice (dynamic array).
   You will need to adjust this implementation to be synchronized in the exact same way you did for the `SafeStack` type.
   This has to be done to make the benchmark between the three implementations fair and comparable.

2. Run the three stack benchmarks using the following command.

   ```console
   go test -v -run none -bench BenchmarkStacks -memprofilerate=1 -benchmem
   ```

   That is we don't run any tests, because we are only interested in the benchmarks, matched by the `-bench BenchmarkStacks` flag.
   The command also enables memory allocation statistics by supplying the `-benchmem` flag, and the `-memprofilerate` controls the fraction of memory allocations that are recorded and reported in the memory profile.
   By passing 1 here means all allocations are reported.

   Attach the benchmark output in your [`benchmark_report.md`](benchmark_report.md) and answer the questions.

3. Run the `CspStack` benchmark separately and write a CPU profile to file:

   ```console
   go test -v -run none -bench BenchmarkStacks/CspStack -cpuprofile=csp-stack.prof
   ```

   Load the CPU profile data with the `pprof` tool.

   ```console
   go tool pprof csp-stack.prof
   ```

   Attach the benchmark and profile output in your [`benchmark_report.md`](benchmark_report.md), and answer the questions related the top ten functions from your CPU profile.

4. Run the `SafeStack` benchmarks separately and write a memory profile to file:

   ```console
   go test -v -run none -bench BenchmarkStacks/SafeStack -memprofile=safe-stack.prof
   ```

   Using the `pprof` tool:

   ```console
   go tool pprof safe-stack.prof
   ```

   Identify the function allocating memory in the `SafeStack` implementation, and list the relevant function to identify the line where the allocations occur.
   Attach the profile output in your [`benchmark_report.md`](benchmark_report.md), and answer the questions related to memory allocations.

5. Install [Graphviz](http://www.graphviz.org/).
   Explore the visualization possibilities offered by `go tool pprof` when analyzing profiling data.
   Use the `pdf` command to produce a call graph:

   ```console
   $ go tool pprof csp-stack.prof
   ...
   (pprof) pdf
   Generating report in profile001.pdf
   (pprof) quit
   ```

   Add the `profile001.pdf` to the `fig/` folder in your group's repository.
   Examine the call graph visualization and answer the questions in the [`benchmark_report.md`](benchmark_report.md).
