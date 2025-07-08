# Introduction

- This repository contains examples of concurrency using golang *Goroutines*.
- *Goroutines* are lightweight threads managed by golang runtime that provide concurrent
processing of functions.

# When to use Goroutines
- As mentioned before, Goroutines provide concurrency for golang, the right question should be, "When should I use concurrency?"
- Concurrency is a strong method of processing that allows code to be executed at the same time.
Example: You have a batch of multiple users payment that needs to be processed through an external API service.
You can write a function that calls this payment API and use Goroutines so all the processing is executed
"at the same time".
- Under the hood Goroutines use execution scheduling, so it just appear that is being executed at the same time, but what
it is really being done is a smart use of processing, where the scheduler can check if a thread is currently not being used (waiting for IO for example, like a HTTP response)
and then put another thread to work so the process is never really stopped, there will always be some piece of code bein ran.