package main

import (
	"fmt"
	"time"
)

func main() {
	// put your code here
	var in string
	fmt.Scan(&in)
	//1986-04-16T05:20:00+06:00

	firstTime, err := time.Parse(time.RFC3339, in)
	if err != nil {
		panic(err)
	}
	//Дана формата:
	//("Mon Jan _2 15:04:05 MST 2006")
	fmt.Println(firstTime.Format(time.UnixDate))
}
