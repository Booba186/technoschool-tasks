package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func worker(id int, jobs <-chan int) {
	for job := range jobs {
		fmt.Printf("Worker %d получил: %d\n", id, job)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <число_воркеров>")
		return
	}

	n, _ := strconv.Atoi(os.Args[1])

	jobs := make(chan int)

	for i := 1; i <= n; i++ {
		go worker(i, jobs)
	}

	counter := 1
	for {
		jobs <- counter
		counter++
		time.Sleep(500 * time.Millisecond)
	}
}
