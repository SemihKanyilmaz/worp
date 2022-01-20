package cmd

import (
	"fmt"
	"time"

	"github.com/semihkanyilmaz/worp/internal/worp"
)

func Execute() {

	w := worp.New()

	w.NewJob("hello-world", 5*time.Second, func() {
		fmt.Println("Hello worp!")
	})
}
