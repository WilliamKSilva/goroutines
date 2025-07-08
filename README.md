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

# Examples

### Page Download
- To run this example you can use `go run page_download`.

- If you want to personalize the test you can pass the flag `-concurrency=true or -concurrency=false` to use Goroutines or not, the default is `true`. You can also pass how much pages should be downloaded using `-pages=100`, the default is `50`.

- On the first example *page_download* the idea is to make multiple GET HTTP request to a random website to download its content and check how much time will take
with and without Goroutines.
Since HTTP requests are IO bound and in this case the response of one download webpage not depend on the others we can spawn a Goroutine to make this operations
concurrently.

- The results *with* Goroutines: running 100 requests our process ran on *~0.29* seconds.
- The results *without* Goroutines: running 100 requests our process ran on *~3.4* seconds. 