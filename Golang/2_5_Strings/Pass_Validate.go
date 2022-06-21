package main

import (
	"fmt"
	"unicode"
)

func main() {
	var a string
	fmt.Scan(&a)
	//fdsghdfgjsdDD1
	ra := []rune(a)
	if len(ra) < 5 {
		fmt.Println("Wrong password")
		return
	}
	for _, let := range a {
		if unicode.IsDigit(let) || unicode.Is(unicode.Latin, let) {
			continue
		} else {
			fmt.Println("Wrong password")
			return
		}
	}
	fmt.Println("Ok")
}
