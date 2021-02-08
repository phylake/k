package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func init() {
	Register("rex", rex)
}

func rex() {
	usage := func() {
		fmt.Fprintln(os.Stderr, "USAGE: rex <regular expression> <index extraction> [index extractions...]")
		os.Exit(2)
	}

	// remove program name and rex
	args := os.Args[2:]
	if len(args) < 2 {
		usage()
	}

	indices := make([]int, 0)
	for _, str := range args[1:] {
		i, err := strconv.Atoi(str)
		if err != nil {
			fmt.Fprintf(os.Stderr, "strconv.Atoi(%s) err=%s\n", str, err.Error())
			usage()
		}

		indices = append(indices, i)
	}

	re, err := regexp.Compile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "regexp.Compile err=%s\n", err.Error())
		os.Exit(2)
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.Stdin.Stat() err=%s\n", err.Error())
		os.Exit(1)
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "os.Stdin is not a valid character device")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		matches := re.FindStringSubmatch(scanner.Text())
		if matches == nil {
			fmt.Fprintln(os.Stderr, "no match")
			os.Exit(1)
		}

		for i, ival := range indices {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(matches[ival])
		}
		fmt.Println("")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner.Err err=%s", err.Error())
		os.Exit(1)
	}
}
