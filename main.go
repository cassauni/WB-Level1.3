package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringSecure(n int) (string, error) {
	b := make([]byte, n)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[num.Int64()]
	}
	return string(b), nil
}

func worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for str := range jobs {
		fmt.Printf("[worker %d] %s\n", id, str)
	}
}

func main() {
	workers := flag.Int("workers", 2, "number of workers")
	flag.Parse()

	if *workers < 1 {
		fmt.Println("workers must be >= 1")
		os.Exit(2)
	}
	jobs := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		str, err := RandStringSecure(*workers)
		if err != nil {
			fmt.Println("rand error:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		select {
		case <-interrupt:
			close(jobs)
			wg.Wait()
			fmt.Println("Done")
			return
		case jobs <- str:
		}
	}
}
