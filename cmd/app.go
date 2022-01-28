package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/semihkanyilmaz/worp/internal/router"
	"github.com/semihkanyilmaz/worp/internal/worp"
)

func Start() {

	w := worp.New()

	w.NewJob("hello-world", 10*time.Second, func() {
		fmt.Println("Hello worp!")
	})
	port := ":1923"
	server := &http.Server{Addr: "127.0.0.1" + port, Handler: router.InitRoutes(w)}

	log.Printf("Http server started on %s", port)

	log.Fatal(server.ListenAndServe())
}
