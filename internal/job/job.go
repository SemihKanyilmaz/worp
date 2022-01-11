package job

import (
	"log"
	"time"
)

type job struct {
	Name         string
	StartSecond  int64
	TimeInterval int64
	Work         func()
}

func NewJob(name string, startSecond, timeInterval int64, work func()) *job {
	return &job{
		Name:         name,
		StartSecond:  startSecond,
		TimeInterval: timeInterval,
		Work:         work,
	}
}

func (j *job) Create() {

	t := time.NewTicker(time.Duration(j.StartSecond) * time.Second)

	go func() {
		for {
			select {
			case <-t.C:
				log.Printf("%s working \n", j.Name)
				go j.Work()
				t.Stop()
				t.Reset(time.Duration(j.TimeInterval) * time.Second)
			}
		}
	}()

}
