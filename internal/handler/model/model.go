package model

import "time"

type GetActiveJobsResponse struct {
	Jobs       []Job `json:"jobs"`
	TotalCount int   `json:"totalCount"`
}

type UpdateDurationJobRequest struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

type Job struct {
	Name      string     `json:"name"`
	Duration  string     `json:"duration"`
	CreatedAt time.Time  `json:"createdAt"`
	LastRunAt *time.Time `json:"lastRunAt"`
	NextRunAt *time.Time `json:"nextRunAt"`
	IsActive  bool       `json:"isActive"`
}
