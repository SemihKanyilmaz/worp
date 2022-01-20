package job_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/semihkanyilmaz/worp/internal/job"
	"github.com/stretchr/testify/assert"
)

func initJob() {
	mu := sync.Mutex{}
	job.InitJob(&mu)
}

func TestCreate_ReturnNil(t *testing.T) {

	initJob()

	job, err := job.NewJob("john-doe", 2, 2, func() {})
	if err == nil {
		t.Error(err)
	}
	err = job.Start()
	if err != nil {
		t.Error(err)
	}

}

//This job name has been already exists!
func TestCreate_ReturnErr(t *testing.T) {

	initJob()
	const mockJobName = "john-doe"

	j, err := job.NewJob(mockJobName, 2, 2, func() {})
	if err != nil && j == nil {
		t.Error(err)
	}

	_, err = job.NewJob(mockJobName, 2, 2, func() {})
	if err == nil {
		t.Errorf("Error must not return nil")
	}

	assert.Equal(t, err.Error(), fmt.Sprintf("%s has been already exists!", mockJobName))

}
