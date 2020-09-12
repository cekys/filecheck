package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	timeStart := time.Now()
	diffs, newFiles, err := compare(DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("diff file:", diffs)
	fmt.Println("new file:", newFiles)
	fmt.Println("work finished, it takes", time.Since(timeStart))
}
