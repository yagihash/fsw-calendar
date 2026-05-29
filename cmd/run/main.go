package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"

	function "github.com/yagihash/fsw-calendar"
)

const (
	exitOK = iota
	exitError
)

func main() {
	os.Exit(run())
}

func run() int {
	data, err := json.Marshal(map[string]string{
		"calendar_id": os.Getenv("CALENDAR_ID"),
		"course":      os.Getenv("COURSE"),
		"class":       os.Getenv("CLASS"),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	if err := function.Register(context.Background(), &pubsub.Message{Data: data}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	return exitOK
}
