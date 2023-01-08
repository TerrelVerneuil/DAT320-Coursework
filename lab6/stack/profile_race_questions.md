# Multiple Choice Questions About Profiling and Data Races

Answer the following questions by editing this file by replacing the `[ ]` for the correct answer with `[x]`.
Only one choice per question is correct.
Selecting more than one choice will result in zero points.
No other changes to the text should be made.

1. What do you call a piece of code that must only be accessed by one process at a time?
    - [ ] a) monitor
    - [ ] b) race condition
    - [ ] c) critical section
    - [ ] d) semaphore

2. Which condition is necessary for a data race?
    - [ ] a) Two or more goroutines tries to read a variable concurrently.
    - [ ] b) Three or more goroutines tries to read a variable concurrently.
    - [ ] c) One or more goroutines tries to read a variable concurrently with a goroutine trying to write to the variable.
    - [ ] d) A single goroutine tries to read a variable followed by a write to the variable.

3. Which command will run the race detector for `TestUnsafeStack`?
    - [ ] a) `go test -race TestUnsafeStack`
    - [ ] b) `go test -race -run TestUnsafeStack`
    - [ ] c) `go test -run TestUnsafeStack`
    - [ ] d) `go test -run -race TestUnsafeStack`

4. Which method(s) in `UnsafeStack` can cause a race condition?
    - [ ] a) `Size()`
    - [ ] b) `Push()`
    - [ ] c) `Pop()`
    - [ ] d) All of the above

5. Which testing flag changes the amount of memory allocations measured?
    - [ ] a) `-benchmem`
    - [ ] b) `-memprofile`
    - [ ] c) `-memprofilerate`
    - [ ] d) `-benchtime`

6. Which `pprof` command will list the top 20 profiles?
    - [ ] a) topN
    - [ ] b) top20
    - [ ] c) top 20
    - [ ] d) 20 top

7. Which `pprof` command will produce a graph of the profile and show it on a web browser?
    - [ ] a) web
    - [ ] b) list
    - [ ] c) svg
    - [ ] d) top

8. Which stack implementation uses the most CPU and memory?
    - [ ] a) SafeStack
    - [ ] b) CspStack
    - [ ] c) SliceStack
