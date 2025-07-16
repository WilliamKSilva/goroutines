package main

import (
	"io"
	"log"
	"net/http"
)

type ServiceCheckResult struct {
	Success    bool
	StatusCode int
	Service    string
	WorkerID   int
}

func serviceCheckWorker(id int, jobs chan string, results chan ServiceCheckResult) {
	for j := range jobs {
		log.Printf("Worker %d started job", id)
		resp, err := http.Get(j)
		if err != nil {
			results <- ServiceCheckResult{
				Service:    j,
				Success:    false,
				StatusCode: resp.StatusCode,
				WorkerID:   id,
			}

			continue
		}

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[ERROR] error parsing the body of the job request for the service: %s, %s", j, err.Error())
			continue
		}

		if resp.StatusCode == http.StatusOK {
			results <- ServiceCheckResult{
				Service:    j,
				Success:    true,
				StatusCode: resp.StatusCode,
				WorkerID:   id,
			}
		} else {
			results <- ServiceCheckResult{
				Service:    j,
				Success:    false,
				StatusCode: resp.StatusCode,
				WorkerID:   id,
			}
		}

		log.Printf("Worker %d finished job", id)
	}
}

func main() {
	servicesToCheck := []string{
		"https://youtube.com",
		"https://facebook.com",
		"https://react.dev",
	}

	numJobs := len(servicesToCheck)
	jobs := make(chan string, numJobs)
	results := make(chan ServiceCheckResult, numJobs)

	for w := 1; w <= 2; w++ {
		go serviceCheckWorker(w, jobs, results)
	}

	for _, s := range servicesToCheck {
		jobs <- s
	}

	close(jobs)

	for r := range results {
		log.Printf("recv service check: %v", r)
	}
}
