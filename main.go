package main

import (
	"fmt"
	"time"

	"../mypkg"
)

func main() {
	timeStart := time.Now()
	err := createJSONTable(DefaultConfig)
	mypkg.ErrorHandler(err, true)
	fmt.Println("work finished, it takes", time.Since(timeStart))
}
