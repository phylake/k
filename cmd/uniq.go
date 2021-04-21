package cmd

import (
	"bufio"
	"fmt"
	"os"
)

func init() {
	Register("uniq", uniq)
}

func uniq() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.Stdin.Stat() err=%s\n", err.Error())
		os.Exit(1)
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "os.Stdin is not a valid character device")
		os.Exit(1)
	}

	msi := make(map[string]interface{})

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msi[scanner.Text()] = nil
	}

	for k := range msi {
		fmt.Println(k)
	}
}
