# Lab 6: Concurrency and Parallelism

| Lab 6: | Concurrency and Parallelism |
| ---------------------    | --------------------- |
| Subject:                 | DAT320 Operating Systems and Systems Programming |
| Deadline:                | **November 1, 2022 23:59** |
| Expected effort:         | 20-30 hours |
| Grading:                 | Pass/fail |
| Submission:              | Group |

## Table of Contents

1. [Introduction](#introduction)
2. [Recommended Reading](#recommended-reading)
3. [Condition Variables](#condition-variables)
4. [Tasks](#tasks)

## Introduction

This lab assignment is divided into three parts and deals with two separate programming tools.
In the first part, you will work with parallel execution using goroutines.
The second part will focus on high-level synchronization techniques and will give an introduction to Go’s built-in data race detector.
You will use two different techniques to ensure synchronization to a shared data structure.
The third part of the lab deals with CPU and memory profiling.
We will analyze different implementations of a simple data structure.

### Recommended Reading

An important resource for this assignment is the [`sync` package](https://golang.org/pkg/sync/) of the standard library.

Further, you may need to read again some of the material from the list of resources from Introduction to Go Programming assignment.
In particular you will want to take a look at chapters about concurrency, goroutines, channels, and synchronization primitives, such as mutex locks and wait groups.
Here are some direct pointers:

* [The Go Programming Language (book)](http://www.gopl.io): Chapters 8 and 9.
* [Collection of Videos about Go](https://github.com/golang/go/wiki/GoTalks), specifically this video about [Concurrency](https://youtu.be/f6kdp27TYZs).
* [Golang Tutorial Series](https://golangbot.com/learn-golang-series/): Sections on Concurrency.

### Condition Variables

As presented in the lectures, condition variables allow us to create sections of code that are thread safe and only execute when some condition is met.
Recall the following two rules of condition variables:

* Always hold the lock when calling `Wait()` and `Signal()`.
* Always recheck the condition after returning from a `Wait()` call.

You can read more about condition variables in Go [here](https://golang.org/pkg/sync/#Cond).

## Tasks

* [Word Count](wordcount/wordcount.md)
* [Concurrent Stack](stack/stack.md)
* [FizzBizz](fizzbizz/fizzbizz.md)
* [Water](water/water.md)
