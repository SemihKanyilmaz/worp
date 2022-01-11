package main

import (
	"fmt"
	"time"

	"github.com/semihkanyilmaz/worp/internal/job"
)

func main() {

	
	job.NewJob("test", 1,5,printMessage).Create()


	time.Sleep(60 * time.Second * 24 * 292)
}

func printMessage(){
	fmt.Println("test")
}

