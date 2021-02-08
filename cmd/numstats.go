package cmd

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func init() {
	Register("numstats", numstats)
}

func numstats() {

	maxLines := math.MaxInt32
	var lines int
	var err error

	args := os.Args[2:]
	if len(args) == 1 {
		maxLines, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "strconv.Atoi() err=%s\n", err.Error())
			os.Exit(2)
		}
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

	nums := make([]float64, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && lines < maxLines {
		i, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\nstrconv.ParseFloat err=%s\n", scanner.Text(), err.Error())
			continue
		}

		lines++
		nums = append(nums, i)
	}

	var (
		sum    float64
		sum2   float64
		max    float64
		min    float64
		count  float64
		avg    float64
		stddev float64
	)
	min = math.MaxFloat64

	for _, ival := range nums {
		sum += ival
		count++

		if ival > max {
			max = ival
		}
		if ival < min {
			min = ival
		}
	}

	if count == 0 {
		min = 0
	}

	fmt.Printf("count: %-40.0f\n", count)
	fmt.Printf("  sum: %-40.0f\n", sum)

	if count > 0 {
		avg = sum / count
		fmt.Printf("  avg: %-40.0f\n", avg)
	}

	for _, ival := range nums {
		sum2 += math.Pow(ival-avg, 2)
	}
	stddev = math.Sqrt(sum2 / count)

	fmt.Printf("    𝜎: %-40.0f\n", stddev)
	fmt.Printf("  max: %-40.0f\n", max)

	sort.Float64s(nums)

	if count >= 10000 {
		p9999 := uint(0.9999 * count)
		fmt.Printf("p9999: %-40.0f\n", nums[p9999])
	}

	if count >= 1000 {
		p999 := uint(0.999 * count)
		fmt.Printf(" p999: %-40.0f\n", nums[p999])
	}

	if count >= 100 {
		p99 := uint(0.99 * count)
		fmt.Printf("  p99: %-40.0f\n", nums[p99])

		p95 := uint(0.95 * count)
		fmt.Printf("  p95: %-40.0f\n", nums[p95])

		p50 := uint(0.50 * count)
		fmt.Printf("  p50: %-40.0f\n", nums[p50])
	}

	fmt.Printf("  min: %-40.0f\n", min)

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner.Err err=%s", err.Error())
		os.Exit(1)
	}
}
