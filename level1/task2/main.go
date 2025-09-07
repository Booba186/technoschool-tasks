package main

import (
	"fmt"
)

func main() {
	numbers := []int{2, 4, 6, 8, 10}
	results := make(chan int, len(numbers))

	for _, num := range numbers {
		go func(n int) {
			square := n * n
			results <- square
		}(num)
	}

	for i := 0; i < len(numbers); i++ {
		result := <-results
		fmt.Println(result)
	}
}
