package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func downloadPage(wg *sync.WaitGroup) {
	if wg != nil {
		log.Println("[WITH GOROUTINES]")
		defer wg.Done()
	} else {
		log.Println("[WITHOUT GOROUTINES]")
	}
	_, err := http.Get("https://www.youtube.com")
	if err != nil {
		log.Printf("[ERROR] error trying to read web page: %v", err)
		os.Exit(1)
	}

	log.Println("Received web page response")
}

func main() {
	var wg sync.WaitGroup
	t := time.Now()

	concurrency := flag.Bool("concurrency", false, "Example with concurrency or not")
	pages := flag.Int("pages", 5, "How many pages should be downloaded")
	flag.Parse()

	if *concurrency {
		wg.Add(*pages)
	}

	for range *pages {
		if *concurrency {
			go downloadPage(&wg)
		} else {
			downloadPage(nil)
		}
	}

	if *concurrency {
		wg.Wait()
	}
	log.Printf("Elapsed time: %f", time.Since(t).Seconds())
}
