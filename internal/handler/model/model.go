package model

import "time"

type GetActiveJobsResponse struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"createdAt"`
	LastRunAt *time.Time `json:"lastRunAt"`
	NextRunAt *time.Time `json:"nextRunAt"`
	IsActive  bool       `json:"isActive"`
}
