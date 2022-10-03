package main

import (
	"context"
	"log"
	"time"

	"github.com/omarahm3/automato/database"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	db, err := database.Connect("", ctx)
	check(err)
	defer cancel()

}

func check(err error) {
	if err == nil {
		return
	}

	log.Fatalln("error occurred", err)
}
