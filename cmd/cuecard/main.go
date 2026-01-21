package main

import (
	"log"

	"github.com/grantcarthew/cuecard/internal/ui"
)

func main() {
	app := ui.New()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
