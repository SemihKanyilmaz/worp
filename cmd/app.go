package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/semihkanyilmaz/worp/internal/router"
	"github.com/semihkanyilmaz/worp/internal/worp"
)

func Execute() {

	w := worp.New()

	w.NewJob("hello-world", 10*time.Second, func() {
		fmt.Println("Hello worp!")
	})

	server := &http.Server{Addr: ":1923", Handler: router.InitRoutes(w)}

	log.Println("Http server started")

	log.Fatal(server.ListenAndServe())
}
