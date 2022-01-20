package httpContext

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Json(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "/api/in-memory", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, r)

	ctx.Json(http.StatusOK, struct {
		Message string `json:"msg"`
	}{
		Message: "test",
	})

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code. Actual %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")

	if contentType != "application/json" {
		t.Errorf("Wrong content type. Actual %s", contentType)
	}

	if strings.Trim(w.Body.String(), " ") != `{"msg":"test"}` {
		t.Errorf("Wrong body. Actual %s", w.Body.String())
	}

}

func Test_InvalidType_Json(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/in-memory", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, r)

	invalidType := make(chan int)

	ctx.Json(http.StatusTeapot, invalidType)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Wrong status code. Actual %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")

	if contentType != "application/json" {
		t.Errorf("Wrong content type. Actual %s", contentType)
	}
}

func Test_Bind(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/in-memory", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, r)

	if err := ctx.Bind(struct{ Message string }{Message: "john-doe"}); err != nil {
		t.Error("Test bind failed")
	}
}

func Test_InvalidType_Bind(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/in-memory", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, r)

	invalidType := []string{}

	if err := ctx.Bind(&invalidType); err == nil {
		t.Error("Test bind failed")
	}

}

func Test_GetQueryParam(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/in-memory?john=doe", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, r)

	if value := ctx.GetQueryParam("john"); value != "doe" {
		t.Errorf("Test GetQueryParam failed. Actual %s ", value)
	}

}

func Test_NoContent(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "/api/in-memory?john=doe", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, r)

	ctx.NoContent(http.StatusCreated)

	contentType := w.Header().Get("Content-Type")

	if contentType != "" {
		t.Errorf("Test NoContent failed. Actual %s", contentType)
	}

	if w.Code != http.StatusCreated {
		t.Errorf("Wrong status code. Actual %d", w.Code)
	}

}
