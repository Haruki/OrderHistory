package main

import (
	"fmt"
	"time"
)

func tmain() {
	fmt.Println("hallo")
	newdate, err := time.Parse("02. Jan. 2006", "24. Jan. 2019")
	if err != nil {
		fmt.Println("fatal error")
	}
	fmt.Println(newdate.Format("2006-01-02"))
}
