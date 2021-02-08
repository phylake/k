package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func init() {
	Register("count", count)
}

func count() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.Stdin.Stat() err=%s\n", err.Error())
		os.Exit(1)
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "os.Stdin is not a valid character device")
		os.Exit(1)
	}

	counts := make(map[string]uint)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		counts[scanner.Text()]++
	}

	var keys []string
	for k := range counts {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		if counts[keys[i]] != counts[keys[j]] {
			return counts[keys[i]] > counts[keys[j]]
		}
		return keys[i] > keys[j]
	})

	for _, k := range keys {
		fmt.Println(counts[k], k)
	}
}
