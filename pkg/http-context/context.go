package httpContext

import (
	"encoding/json"
	"log"
	"net/http"
)

type Context interface {
	Json(statusCode int, i interface{})
	Bind(model interface{}) error
	GetQueryParam(param string) string
	NoContent(statusCode int)
}

type handlerContext struct {
	r *http.Request
	w http.ResponseWriter
}

func NewContext(w http.ResponseWriter, r *http.Request) *handlerContext {
	return &handlerContext{
		r: r,
		w: w,
	}
}

func (c *handlerContext) Json(statusCode int, i interface{}) {
	c.w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(i)
	if err != nil {
		go log.Println(err.Error())
		c.w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.w.WriteHeader(statusCode)
	c.w.Write(data)
}

func (c *handlerContext) Bind(model interface{}) error {

	if err := json.NewDecoder(c.r.Body).Decode(&model); err != nil {
		go log.Printf("Body parse error. Error : %s", err.Error())
		return err
	}
	return nil
}

func (c *handlerContext) GetQueryParam(param string) string {
	return c.r.URL.Query().Get(param)
}

func (c *handlerContext) NoContent(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *handlerContext) Recover() {
	if err := recover(); err != nil {
		go log.Printf("Actually this is a panic message : %v", err)
		c.w.WriteHeader(http.StatusInternalServerError)
	}
}
