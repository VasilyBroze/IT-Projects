package main

import "fmt"

func main() {
	var a, b int
	fmt.Scan(&a)
	// 5
	slice := make([]int, a)
	for i := 0; i < a; i++ {
		fmt.Scan(&slice[i])
		// 1 2 3 -1 -4
		if slice[i] > 0 {
			b = b + 1
		}

	}
	fmt.Print(b)
}
