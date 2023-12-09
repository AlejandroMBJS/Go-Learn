package main

import (
	"fmt"
)

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	
	var num int
	fmt.Print("Enter a number: ")
	fmt.Scanln(&num)

	result := factorial(num)
	fmt.Printf("Factorial of %d is %d\n", num, result)
}