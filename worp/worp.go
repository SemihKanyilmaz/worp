package worp

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Worp interface {
	CreateJob(name string, durat time.Duration, task func()) (*Job, error)
	DeleteJob(name string) error
	Start(name string) error
	PauseJob(name string) error
	GetActiveJobs() []Job
	UpdateDuration(name string, durat time.Duration) error
	UpdateNextRunAt(name string, nextRunAt time.Time) error
}

type worp struct {
	jobs map[string]*Job
	mu   *sync.RWMutex
}

func New() *worp {
	return &worp{
		mu:   &sync.RWMutex{},
		jobs: make(map[string]*Job),
	}
}

type Job struct {
	Name      string
	ticker    *time.Ticker
	LastRunAt *time.Time
	NextRunAt *time.Time
	Durat     time.Duration
	CreatedAt time.Time
	IsActive  bool
	task      func()
}

func (w *worp) CreateJob(name string, durat time.Duration, task func()) (*Job, error) {

	if j, _ := w.getJob(name); j != nil {
		return nil, fmt.Errorf("%s has been already exists", name)
	}

	j := &Job{
		Name:      name,
		ticker:    time.NewTicker(durat),
		task:      task,
		CreatedAt: time.Now(),
		Durat:     durat,
	}

	w.mu.Lock()
	w.jobs[name] = j
	w.mu.Unlock()

	return j, nil
}

func (w *worp) Start(name string) error {

	j, err := w.getJob(name)
	if err != nil {
		return err
	}

	if !j.IsActive {
		j.ticker.Reset(j.Durat)
		j.IsActive = true
	}

	t := time.Now()
	w.mu.Lock()
	j.LastRunAt = &t
	nexRunAt := t.Add(j.Durat)
	j.NextRunAt = &nexRunAt
	w.mu.Unlock()
	log.Printf("%s is working \n", j.Name)

	go func() {
		select {
		case <-j.ticker.C:
			w.mu.Lock()
			j.ticker.Stop()
			j.ticker.Reset(j.Durat)

			now := time.Now()
			j.LastRunAt = &now
			runAt := now.Add(j.Durat)
			j.NextRunAt = &runAt
			w.mu.Unlock()
			j.task()
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

	log.Printf("%s was successfully deleted", j.Name)
	return nil
}

func (w *worp) PauseJob(name string) error {

	j, err := w.getJob(name)
	if err != nil {
		return err
	}
	w.mu.Lock()
	j.IsActive = false

	j.ticker.Stop()
	w.mu.Unlock()
	log.Printf("%s was successfully paused", j.Name)

	return nil
}

func (w *worp) GetActiveJobs() []Job {

	items := make([]Job, 0)

	w.mu.Lock()
	defer w.mu.Unlock()

	for _, job := range w.jobs {
		items = append(items, *job)
	}

	return items
}

func (w *worp) UpdateDuration(name string, durat time.Duration) error {

	j, err := w.getJob(name)
	if err != nil {
		return err
	}

	w.mu.Lock()
	j.Durat = durat

	if j.IsActive {
		j.ticker.Stop()
		j.ticker.Reset(j.Durat)

		now := time.Now()
		j.LastRunAt = &now
		runAt := now.Add(j.Durat)
		j.NextRunAt = &runAt
	}
	w.mu.Unlock()
	return nil
}

func (w *worp) UpdateNextRunAt(name string, nextRunAt time.Time) error {

	if nextRunAt.Before(time.Now().Add(3 * time.Hour)) {
		return errors.New("nextRunAt must be greater than now")
	}

	j, err := w.getJob(name)
	if err != nil {
		return err
	}

	oldDurat := j.Durat

	w.mu.Lock()
	j.ticker.Stop()
	j.ticker.Reset(nextRunAt.Sub(time.Now().Add(3 * time.Hour)))
	j.ticker = time.NewTicker(oldDurat)

	j.NextRunAt = &nextRunAt
	w.mu.Unlock()
	return nil
}

func (w *worp) getJob(name string) (*Job, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	j, exist := w.jobs[name]
	if !exist {
		return nil, fmt.Errorf("%s not found in jobs", name)
	}
	return j, nil
}
