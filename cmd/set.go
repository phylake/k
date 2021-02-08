package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func init() {
	Register("set", set)
}

func set() {
	usage := func() {
		fmt.Fprintln(os.Stderr, "USAGE: set [diff|int|union] <path to set file 1> <path to set file 2>")
	}

	// remove program name and set
	args := os.Args[2:]
	if len(args) != 3 {
		usage()
		os.Exit(1)
	}

	op := args[0]
	path1 := args[1]
	path2 := args[2]

	f1, err := os.Open(path1)
	if err != nil {
		fmt.Println("os.Open", "err", err)
		os.Exit(1)
	}
	defer f1.Close()

	f2, err := os.Open(path2)
	if err != nil {
		fmt.Println("os.Open", "err", err)
		os.Exit(1)
	}
	defer f2.Close()

	s1 := make(map[string]interface{})
	s2 := make(map[string]interface{})

	scanner1 := bufio.NewScanner(f1)
	for scanner1.Scan() {
		s1[scanner1.Text()] = nil
	}

	scanner2 := bufio.NewScanner(f2)
	for scanner2.Scan() {
		s2[scanner2.Text()] = nil
	}

	if err := scanner1.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner1.Err err=%s", err.Error())
		os.Exit(1)
	}

	if err := scanner2.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner2.Err err=%s", err.Error())
		os.Exit(1)
	}

	s3 := []string{}
	switch op {

	case "diff":
		for k := range s1 {
			if _, exists := s2[k]; !exists {
				s3 = append(s3, k)
			}
		}

	case "int":
		for k := range s1 {
			if _, exists := s2[k]; exists {
				s3 = append(s3, k)
			}
		}

	case "union":
		uniqueList := make(map[string]interface{})

		for k := range s1 {
			uniqueList[k] = nil
		}

		for k := range s2 {
			uniqueList[k] = nil
		}

		for k := range uniqueList {
			s3 = append(s3, k)
		}
	default:
		usage()
		os.Exit(1)
	}

	sort.Strings(s3)
	for _, k := range s3 {
		fmt.Println(k)
	}
}
