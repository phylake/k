// Copyright 2021 Brandon Cook
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
// may be used to endorse or promote products derived from this software without
// specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/phylake/k/cmd"
)

func main() {

	if f, exists := cmd.GetCommand(os.Args[1]); exists {
		f()
	} else {
		args := cmdArgs()
		printArgs(args)
		cmd := exec.Command("kubectl", args...)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			os.Exit(1)
		}
	}
}

func printArgs(args []string) {
	iargs := make([]interface{}, len(args)+1)
	iargs[0] = "kubectl"
	for i := 0; i < len(args); i++ {
		iargs[i+1] = args[i]
	}
	fmt.Fprintln(os.Stderr, iargs...)
}

func cmdArgs() []string {
	args := os.Args[1:]
	argsMap := make(map[string]int)
	for i, arg := range args {
		argsMap[arg] = i
	}

	ns := os.Getenv("NS")
	if ns != "" {
		haveNs := false
		for _, arg := range args {
			arg = strings.TrimSpace(arg)
			if arg == "-n" || arg == "--namespace" || arg == "--all-namespaces" {
				haveNs = true
				break
			}
		}
		if !haveNs {
			args = append([]string{"-n", ns}, args...)
		}
	}

	return args
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
