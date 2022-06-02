package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func init() {
	Register("time", timeF)
}

func timeF() {
	usage := func(err error) {
		fmt.Fprintf(os.Stderr, "USAGE: time <duration string> [%s]\n",
			time.Now().Format(time.RFC3339))
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// remove program name and time
	args := os.Args[2:]
	if len(args) != 1 && len(args) != 2 {
		usage(errors.New(""))
	}

	dur, err := time.ParseDuration(args[0])
	if err != nil {
		usage(err)
	}

	now := time.Now().UTC()
	if len(args) == 2 {
		t, err := time.Parse(time.RFC3339, args[1])
		if err != nil {
			usage(err)
		}
		now = t
	}

	fmt.Println(now.Add(dur).Format(time.RFC3339))
}
