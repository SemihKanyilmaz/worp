package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/semihkanyilmaz/worp/internal/handler/model"
	"github.com/semihkanyilmaz/worp/internal/worp"
	httpContext "github.com/semihkanyilmaz/worp/pkg/http-context"
)

type handler struct {
	worp worp.Worp
}

func NewHandler(worp worp.Worp) *handler {
	return &handler{worp: worp}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c := httpContext.NewContext(w, r)
	defer c.Recover()

	switch r.Method {
	case http.MethodGet:
		h.getActiveJobs(c)
	case http.MethodPost:
		h.restartJob(c)
	case http.MethodDelete:
		h.deleteJob(c)
	case http.MethodPut:
		h.updateTimeJob(c)
	default:
		c.Json(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

}

func (h *handler) PauseJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	c := httpContext.NewContext(w, r)
	name := strings.TrimPrefix(r.URL.Path, "/pause/")
	if strings.TrimSpace(name) == "" {
		c.Json(http.StatusBadRequest, "name cannot be empty")
		return
	}

	if err := h.worp.PauseJob(name); err != nil {

		c.Json(http.StatusNotFound, err.Error())
		return
	}

	c.NoContent(http.StatusNoContent)
}

func (h *handler) getActiveJobs(c httpContext.Context) {

	jobs := h.worp.GetActiveJobs()
	var res model.GetActiveJobsResponse
	length := len(jobs)
	res.Jobs = make([]model.Job, length)
	res.TotalCount = length

	for i, job := range jobs {
		res.Jobs[i].CreatedAt = job.CreatedAt
		res.Jobs[i].Duration = job.Durat.String()
		res.Jobs[i].IsActive = job.IsActive
		res.Jobs[i].LastRunAt = job.LastRunAt
		res.Jobs[i].NextRunAt = job.NextRunAt
		res.Jobs[i].Name = job.Name
	}

	c.Json(http.StatusOK, res)

}

func (h *handler) updateTimeJob(c httpContext.Context) {

	req := new(model.UpdateTimeRequest)
	if err := c.Bind(req); err != nil {
		c.Json(http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		c.Json(http.StatusBadRequest, "name cannot be empty")
		return
	}

	if req.Duration != nil {
		if err := h.worp.UpdateDuration(req.Name, *req.Duration); err != nil {
			c.Json(http.StatusBadRequest, err.Error())
			return
		}
	}

	if req.NextRunAt != nil {

		next, err := time.Parse(time.RFC3339, *req.NextRunAt)
		if err != nil {
			c.Json(http.StatusBadRequest, err.Error())
			return
		}

		if err := h.worp.UpdateNextRunAt(req.Name, next); err != nil {
			c.Json(http.StatusNotFound, err.Error())
			return
		}
	}

	c.NoContent(http.StatusNoContent)
}

func (h *handler) deleteJob(c httpContext.Context) {

	name := c.GetQueryParam("name")
	if strings.TrimSpace(name) == "" {
		c.Json(http.StatusBadRequest, "name cannot be empty")
		return
	}

	if err := h.worp.DeleteJob(name); err != nil {
		c.Json(http.StatusNotFound, err.Error())
		return
	}

	c.NoContent(http.StatusNoContent)
}

func (h *handler) restartJob(c httpContext.Context) {

	name := c.GetQueryParam("name")
	if strings.TrimSpace(name) == "" {
		c.Json(http.StatusBadRequest, "Name cannot be empty")
		return
	}
	err := h.worp.Start(name)
	if err != nil {
		c.Json(http.StatusNotFound, err.Error())
		return
	}

	c.Json(http.StatusOK, "Job has been succesffully restarted!")
}
