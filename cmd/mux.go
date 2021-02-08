package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func init() {
	Register("mux", mux)
}

func mux() {
	if os.Getenv("SELECTOR") == "" {
		fmt.Println("SELECTOR must be set in environment")
		os.Exit(2)
	}

	if os.Getenv("NS") == "" {
		fmt.Println("NS must be set in environment")
		os.Exit(2)
	}

	// remove program name and mux
	args := os.Args[2:]
	switch args[0] {
	case "exec":
		muxExec(args[1:])
	case "logs":
		fmt.Println("TODO")
		os.Exit(3)
	default:
		fmt.Println("USAGE: mux [exec|logs]")
		os.Exit(2)
	}
}

func muxExec(execArgs []string) {
	selector := os.Getenv("SELECTOR")

	var buf bytes.Buffer

	cmd := exec.Command("kubectl",
		"-n", os.Getenv("NS"),
		"get", "pods",
		"--selector", selector,
		"-o=jsonpath='{.items[*].metadata.name}'")

	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	pods := strings.Split(strings.Trim(buf.String(), "'"), " ")
	for _, pod := range pods {
		cmdArgs := []string{"-n", os.Getenv("NS"), "exec", pod}
		cmdArgs = append(cmdArgs, execArgs...)

		printArgs := []interface{}{"kubectl"}
		for _, arg := range cmdArgs {
			printArgs = append(printArgs, arg)
		}
		fmt.Fprintln(os.Stderr, printArgs...)

		cmd := exec.Command("kubectl", cmdArgs...)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func remove(args []string, delArg string) []string {
	for i, arg := range args {
		if arg == delArg {
			args = append(args[:i], args[i+1:]...)
		}
	}
	return args
}

func insertAfter(args []string, index int, insertArgs ...string) []string {
	for i := range args {
		if index == i {
			args2 := append(insertArgs, args[i:]...)
			args = append(args[:i], args2...)
			break
		}

	}
	return args
}
