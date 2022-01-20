package job

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var jobs = map[string]job{}

var mu *sync.Mutex

func InitJob(mutex *sync.Mutex) {
	mu = mutex
}

type job struct {
	Name         string
	TimeInterval int64
	ticker       *time.Ticker
	isActive     bool
	Work         func()
}

func NewJob(name string, startSecond, timeInterval int64, work func()) (*job, error) {
	mu.Lock()
	if _, exist := jobs[name]; exist {
		return nil, fmt.Errorf("%s has been already exists!", name)
	}
	mu.Unlock()
	j := job{
		Name:         name,
		ticker:       time.NewTicker(time.Duration(startSecond) * time.Second),
		Work:         work,
		TimeInterval: timeInterval,
		isActive:     true,
	}
	mu.Lock()
	jobs[name] = j
	mu.Unlock()

	return &j, nil
}

func (j *job) Start() error {

	checkMutexNil()

	if _, exist := jobs[j.Name]; !exist {
		mu.Lock()
		jobs[j.Name] = *j
		mu.Unlock()
	}

	if !j.isActive {
		j.ticker.Reset(time.Duration(j.TimeInterval) * time.Second)
		j.isActive = true
	}

	log.Printf("Job named %s working \n", j.Name)

	go func() {
		for {
			select {
			case <-j.ticker.C:
				go j.Work()
				j.ticker.Stop()
				j.ticker.Reset(time.Duration(j.TimeInterval) * time.Second)
			}
		}
	}()

	return nil
}

func (j *job) DeleteJob() error {

	checkMutexNil()

	if err := isJobExists(j.Name); err != nil {
		return err
	}

	mu.Lock()
	delete(jobs, j.Name)
	mu.Unlock()

	j.isActive = false

	j.ticker.Stop()

	log.Printf("Job named %s was successfully deleted", j.Name)
	return nil
}

func (j *job) PauseJob() error {
	checkMutexNil()
	if err := isJobExists(j.Name); err != nil {
		return err
	}

	j.isActive = false

	j.ticker.Stop()

	log.Printf("Job named %s was successfully paused!", j.Name)

	return nil
}

func GetActiveJobs() []job {
	checkMutexNil()

	items := make([]job, 0)

	mu.Lock()

	for _, job := range jobs {
		if job.isActive {

			items = append(items, job)
		}
	}
	mu.Unlock()

	return items
}

func isJobExists(name string) error {
	mu.Lock()
	if _, exist := jobs[name]; !exist {
		return fmt.Errorf("%s not found!", name)
	}
	mu.Unlock()
	return nil
}

func checkMutexNil() {
	if mu == nil {
		log.Fatal("Before using run NewJob function, initiliaze InitJob function")
	}
}
