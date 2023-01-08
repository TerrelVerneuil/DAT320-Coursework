# Task: Parallel Word Count

Parallel execution is often quite simple in Go.
In this part you will implement a function to perform parallel word count.
We have provided you with a large text file (`mobydick.txt`) and three simple functions in the file `wc.go` you should use:

* The function `loadMoby()` loads the file and returns it as a `[]byte`.
* The function `wordCount()` counts the words in a `[]byte`.
* And the function `shardSlice()` splits the provided `[]byte` into sub-slices that can be counted separately.

1. Your first task is to implement a function called `parallelWordCount()` with the following signature:

    ```go
    func parallelWordCount(input []byte) (words int)
    ```

    This function *must count* the words in `mobydick.txt` using multiple goroutines; typically as many as there are CPU cores on your machine.
    The parallel version must support different number of goroutines.
    You can use the provided functions as you see fit.
    The `TestParallelWordCount` must pass, meaning that it should return the same number of words as the sequential version of word count (the `wordCount()` function).

    To run the local tests:

    ```console
    go test -v -run TestParallelWordCount
    ```

2. Perform benchmarks using the provided benchmark tests. Run as follows:

   ```console
   go test -v -run XX -bench BenchmarkWordCount
   ```

   Your parallel implementation must perform better than the provided sequential implementation.
   The TAs will check this during approval.
