package main

import (
	"time"

	"github.com/semihkanyilmaz/worp/internal/job"
)

func main() {

	job.NewJob("test", 1, 5, func() {}).Create()

	time.Sleep(60 * time.Second * 24 * 292)
}
