package main

import (
	"PRService/internal/app"
	"context"
	"log"
)

func main() {
	a := app.App{}
	err := a.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
