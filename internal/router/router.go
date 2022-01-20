package router

import (
	"net/http"

	"github.com/semihkanyilmaz/worp/internal/handler"
	"github.com/semihkanyilmaz/worp/internal/worp"
)

func InitRoutes(worp worp.Worp) *http.ServeMux {

	mux := http.NewServeMux()
	handler := handler.NewHandler(worp)
	mux.Handle("/", handler)
	mux.HandleFunc("/pause", handler.PauseJob)
	return mux
}
