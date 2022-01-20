package worp

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Worp interface {
	NewJob(name string, durat time.Duration, work func()) (*job, error)
	DeleteJob(name string) error
	Start(name string) error
	PauseJob(name string) error
	GetActiveJobs() []job
}

type worp struct {
	jobs map[string]*job
	mu   *sync.Mutex
}

func New() *worp {
	return &worp{
		mu:   &sync.Mutex{},
		jobs: make(map[string]*job),
	}
}

type job struct {
	Name     string
	ticker   *time.Ticker
	durat    time.Duration
	isActive bool
	Work     func()
}

func (w *worp) NewJob(name string, durat time.Duration, work func()) (*job, error) {

	if j, _ := w.getJob(name); j != nil {
		return nil, fmt.Errorf("%s has been already exists!", name)
	}

	j := job{
		Name:     name,
		ticker:   time.NewTicker(durat),
		Work:     work,
		durat:    durat,
		isActive: true,
	}
	w.mu.Lock()
	w.jobs[name] = &j
	w.mu.Unlock()

	return &j, nil
}

func (w *worp) Start(name string) error {

	j, err := w.getJob(name)
	if err != nil {
		return err
	}

	if !j.isActive {
		j.ticker.Reset(j.durat)
		j.isActive = true
	}

	log.Printf("Job named %s working \n", j.Name)

	go func() {
		for {
			select {
			case <-j.ticker.C:
				go j.Work()
				j.ticker.Stop()
				j.ticker.Reset(j.durat)
			}
		}
	}()

	return nil
}

func (w *worp) DeleteJob(name string) error {

	j, err := w.getJob(name)
	if err != nil {
		return err
	}

	w.mu.Lock()
	delete(w.jobs, j.Name)
	w.mu.Unlock()

	j.isActive = false

	j.ticker.Stop()

	log.Printf("Job named %s was successfully deleted", j.Name)
	return nil
}

func (w *worp) PauseJob(name string) error {

	j, err := w.getJob(name)
	if err != nil {
		return err
	}

	j.isActive = false

	j.ticker.Stop()

	log.Printf("Job named %s was successfully paused!", j.Name)

	return nil
}

func (w *worp) GetActiveJobs() []job {

	items := make([]job, 0)

	w.mu.Lock()
	defer w.mu.Unlock()

	for _, job := range w.jobs {
		if job.isActive {

			items = append(items, *job)
		}
	}

	return items
}

func (w *worp) getJob(name string) (*job, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	j, exist := w.jobs[name]
	if !exist {
		return nil, fmt.Errorf("%s not found in jobs", name)
	}
	return j, nil
}
