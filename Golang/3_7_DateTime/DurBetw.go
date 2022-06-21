package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	a, err := bufio.NewReader(os.Stdin).ReadString('\n')
	//13.03.2018 14:00:15,12.03.2018 14:00:15
	strings.Trim(a, "\n")
	fmt.Println(a)
	if err != nil {
		panic(err)
	}

	s := strings.Split(a, ",")
	s[1] = strings.TrimSpace(s[1])
	date1, err := time.Parse("02.01.2006 15:04:05", s[0])
	if err != nil {
		panic(err)
	}

	date2, err := time.Parse("02.01.2006 15:04:05", s[1])
	if err != nil {
		panic(err)
	}

	if date1.After(date2) {
		otv := date1.Sub(date2)
		fmt.Println(otv)
	} else {
		otv := date2.Sub(date1)
		fmt.Println(otv)
	}

}
