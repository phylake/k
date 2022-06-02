k tool
======

`kubectl` wrapper and utilities

- [Basic usage](#basic-usage)
- [Utilities](#utilities)
- [Extensibility](#extensibility)

Basic usage
-----------

Type `k` instead of `kubectl` and save 6 keystrokes.

The command to be run will be printed to stderr.

```bash
$ k get nodes
kubectl get nodes
NAME                                         STATUS   ROLES    AGE   VERSION
ip-10-10-1-204.us-west-2.compute.internal    Ready    <none>   60m   v1.18.14
```

Avoid typing `-n <ns>` over and over by exporting `NS`

```bash
$ export NS=some-long-ns-name
$ k get pods
kubectl -n some-long-ns-name get pods
NAME                     READY   STATUS    RESTARTS   AGE
some-pod                 1/1     Running   0          10m
````

Utilities
--------

- [mux](#mux)
- [rex](#rex)
- [numstats](#numstats)
- [uniq](#uniq)
- [count](#count)
- [set](#set)
- [time](#time)

### mux

Multiplex `exec` command output across Pods onto stdout

Usage: `k mux exec <command>`

Useful for gathering output from Pods by selector. `NS` is still set from above.

```bash
$ SELECTOR=app=envoy k mux exec echo hi
kubectl -n some-long-ns-name exec envoy-65f74d8949-5sk9b echo hi
hi
kubectl -n some-long-ns-name exec envoy-65f74d8949-5sk9b echo hi
hi
kubectl -n some-long-ns-name exec envoy-65f74d8949-5sk9b echo hi
hi
```

The command is sent to stderr so output can captured in the clipboard or a file

```bash
$ SELECTOR=app=envoy k mux exec echo hi | pbcopy
kubectl -n some-long-ns-name exec envoy-65f74d8949-5sk9b echo hi
kubectl -n some-long-ns-name exec envoy-65f74d8949-5sk9b echo hi
kubectl -n some-long-ns-name exec envoy-65f74d8949-5sk9b echo hi

$ pbpaste
hi
hi
hi
```

### rex

Manipulate stdin with a [regular expression](https://regex101.com/) and output capture groups

Usage: `k rex <regular expression> <capture group 1> [capture group 2] [capture group 3]...`

```bash
$ echo 'abc 123\ndef 456'
abc 123
def 456

$ echo 'abc 123\ndef 456' | k rex '\w+' 0
abc
def

$ echo 'abc 123\ndef 456' | k rex '\d+' 0
123
456

$ echo 'abc 123\ndef 456' | k rex '(\w+) (\d+)' 2 1
123 abc
456 def
```

Top pods example

```bash
$ k top pods
kubectl -n some-long-ns-name top pods
NAME                     CPU(cores)   MEMORY(bytes)
envoy-65f74d8949-5sk9b   33m          75Mi
envoy-65f74d8949-87h54   38m          77Mi
envoy-65f74d8949-rrx9z   36m          75Mi

$ k top pods | grep envoy | k rex '((\w|-)+) +(\d+)' 3 1 | sort -rn
kubectl -n some-long-ns-name top pods
38 envoy-65f74d8949-87h54
36 envoy-65f74d8949-rrx9z
33 envoy-65f74d8949-5sk9b
```

### numstats

Run basic statistics on stdin

Usage: `k numstats` (assumes numerical input on stdin)

```bash
$ echo '1\n2\n3'
1
2
3

$ echo '1\n2\n3' | k numstats
count: 3
  sum: 6
  avg: 2
    𝜎: 1
  max: 3
  min: 1
```

Look for even load distribution on Pods

```bash
$ k top pods | grep envoy | k rex '(\w|-)+ +(\d+)' 2 | k numstats
kubectl -n some-long-ns-name top pods
count: 151
  sum: 95139
  avg: 630
    𝜎: 207
  max: 1131
  p99: 1130
  p95: 963
  p50: 626
  min: 523
```

### uniq

Print unique values from stdin

```bash
$ echo 'one\ntwo\ntwo\nthree\nthree\nthree'
one
two
two
three
three
three
```

```bash
$ echo 'one\ntwo\ntwo\nthree\nthree\nthree' | k uniq
one
two
three
```

### count

Similar to `uniq` but count stdin and sort in descending order

Usage: `k count`

```bash
$ echo 'one\ntwo\ntwo\nthree\nthree\nthree'
one
two
two
three
three
three

$ echo 'one\ntwo\ntwo\nthree\nthree\nthree' | k count
3 three
2 two
1 one
```

### set

Set operations (intersection, difference, union) on 2 files

Usage: `k set <int|diff|union> <file 1> <file 2>`

Given file 1 at `~/set1` with

```
a
b
c
```

and file 2 at `~/set2` with

```
c
d
e
```

```bash
$ k set int ~/set1 ~/set2
c

$ k set diff ~/set1 ~/set2
a
b

$ k set union ~/set1 ~/set2
a
b
c
d
e
```

### time

Time deltas from now or a RFC 3339 provided time

No arguments prints usage and the current local time which is convenient to then
copy and alter as the second argument (see below)

```bash
$ k time
USAGE: time <duration string> [2022-06-02T15:04:37-06:00]
```

Replace the `-06:00` with `Z` to get UTC

If the current date and time was `2022-01-01T00:00:00Z` then this would add an
hour to it

```bash
$ k time 1h
2022-01-01T01:00:00Z
```

If the current data and time is no longer `2022-01-01T00:00:00Z` but you want to
calculate time deltas from then you can pass it as the second argument. Negative
durations can also be used

```bash
$ k time -1h 2022-01-01T00:00:00Z
2021-12-31T23:00:00Z
```

The available units for the duration are described in [`time.ParseDuration`](https://pkg.go.dev/time@go1.18.2#ParseDuration)

Notably `d` for day is missing which is often inconvenient

Extensibility
-------------

Drop new command files into `cmd/` and register them like this

```go
func init() {
	Register("mycmd", mycommand)
}

func mycommand() {
}
```

Or import `"github.com/phylake/k/cmd"` and call `cmd.Register()`
