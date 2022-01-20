package handler

import (
	"net/http"
	"strings"

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
	default:
		c.Json(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

}

func (h *handler) PauseJob(w http.ResponseWriter, r *http.Request) {
	c := httpContext.NewContext(w, r)
	name := c.GetQueryParam("name")
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

	items := make([]model.Job, len(jobs), len(jobs))

	for i, job := range jobs {
		items[i].Name = job.Name
	}

	c.Json(http.StatusOK, model.GetActiveJobsResponse{
		Jobs: items,
	})

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

	err := h.worp.Start(name)
	if err != nil {
		c.Json(http.StatusNotFound, err.Error())
		return
	}

	c.Json(http.StatusOK, "Job has been succesffully restarted!")
}
