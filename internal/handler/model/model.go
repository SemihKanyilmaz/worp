package model

type GetActiveJobsResponse struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	Name string `json:"name"`
}
