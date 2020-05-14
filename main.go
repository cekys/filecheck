package main

import (
	"fmt"
	"time"

	"../mypkg"
)

func main() {
	timeStart := time.Now()
	diffs, newFiles, err := compare(DefaultConfig)
	//err := createJSONTable(DefaultConfig, false)
	mypkg.ErrorHandler(err, true)
	fmt.Println("diff file:", diffs)
	fmt.Println("new file:", newFiles)
	fmt.Println("work finished, it takes", time.Since(timeStart))
}
