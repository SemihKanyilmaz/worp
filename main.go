package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/semihkanyilmaz/worp/internal/job"
)

func init() {
	mu := sync.Mutex{}

	job.InitJob(&mu)
}

func main() {

	j, err := job.NewJob("test", 1, 1, printMessage)
	if err != nil {
		log.Fatal(err)
	}

	err = j.Start()
	if err != nil {
		log.Println(err)
	}

	j.Start()

	select {}
}

func printMessage() {

	fmt.Println("Hello worp")
}
