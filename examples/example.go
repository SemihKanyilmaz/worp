package main

import (
	"fmt"
	"log"
	"time"

	"github.com/semihkanyilmaz/worp/worp"
)

func main() {
	w := worp.New()

	//A job is defined to run every 3 seconds
	j, err := w.CreateJob("hello-worp", 3*time.Second, doSomething("Hello worp"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	//The job will start
	w.Start(j.Name)

	time.Sleep(7 * time.Second)

	//The job that before initialized, now it will run every 1 second
	if err := w.UpdateDuration(j.Name, time.Second); err != nil {
		log.Fatal(err.Error())
	}

	//The job that before initialized, it will start 3 hours later
	w.UpdateNextRunAt(j.Name, time.Now().Add(3*time.Hour))

	//The job will pause
	w.PauseJob(j.Name)

	//The job will delete
	w.DeleteJob(j.Name)

	//It returns all initialized jobs
	w.GetActiveJobs()

}

func doSomething(msg string) func() {
	return func() {
		fmt.Println(msg)
	}
}
