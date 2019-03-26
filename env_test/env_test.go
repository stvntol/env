package env_test

import (
	. "github.com/stvntol/env"
	"path"
	"strings"
	"testing"
)

func splitPath(p string, depth int) (head, tail string) {
	p = path.Clean("/" + p)

	if depth < 0 {
		return "", p
	}

	parts := strings.Split(p, "/")[1:]
	if depth >= len(parts) {
		return "", "/"
	}

	head = parts[depth]
	tail = "/" + strings.Join(parts[depth+1:], "/")

	return head, tail
}

func testSplitShiftPath(fn func(string, int) (string, string), t *testing.T) {

	p := "/a/b/c/d"
	p2 := "a/b/c/d/"

	type result struct {
		head string
		tail string
	}

	cases := []struct {
		path  string
		depth int
		r     result
	}{
		{p, -1, result{"", "/a/b/c/d"}},
		{p, 0, result{"a", "/b/c/d"}},
		{p, 1, result{"b", "/c/d"}},
		{p, 2, result{"c", "/d"}},
		{p, 3, result{"d", "/"}},
		{p, 4, result{"", "/"}},
		{p, 5, result{"", "/"}},
		{p2, 2, result{"c", "/d"}},
	}

	for _, trial := range cases {
		head, tail := fn(trial.path, trial.depth)
		if head != trial.r.head || tail != trial.r.tail {
			t.Errorf("Got head: %q, tail: %q wanted head: %q, tail: %q",
				head, tail, trial.r.head, trial.r.tail)
		}
	}

}

var alpha string
var beta string

func benchmarkSplitPath(d int, b *testing.B) {
	// run the  function b.N times
	p := "/a/b/c/d"
	var hd string
	var tl string
	for n := 0; n < b.N; n++ {
		hd, tl = splitPath(p, d)
	}
	alpha = hd
	beta = tl
}

func benchmarkShiftPath(d int, b *testing.B) {
	// run the  function b.N times
	p := "/a/b/c/d"
	var hd string
	var tl string
	for n := 0; n < b.N; n++ {
		hd, tl = ShiftPath(p, d)
	}
	alpha = hd
	beta = tl

}

func TestShiftPath(t *testing.T) { testSplitShiftPath(ShiftPath, t) }
func TestSplitPath(t *testing.T) { testSplitShiftPath(splitPath, t) }

func BenchmarkShiftPath0(b *testing.B) { benchmarkShiftPath(0, b) }
func BenchmarkShiftPath1(b *testing.B) { benchmarkShiftPath(1, b) }
func BenchmarkShiftPath2(b *testing.B) { benchmarkShiftPath(2, b) }
func BenchmarkShiftPath3(b *testing.B) { benchmarkShiftPath(3, b) }
func BenchmarkShiftPath4(b *testing.B) { benchmarkShiftPath(4, b) }

func BenchmarkSplitPath0(b *testing.B) { benchmarkSplitPath(0, b) }
func BenchmarkSplitPath1(b *testing.B) { benchmarkSplitPath(1, b) }
func BenchmarkSplitPath2(b *testing.B) { benchmarkSplitPath(2, b) }
func BenchmarkSplitPath3(b *testing.B) { benchmarkSplitPath(3, b) }
func BenchmarkSplitPath4(b *testing.B) { benchmarkSplitPath(4, b) }

func TestNewEnv(t *testing.T) {
	e := NewEnv(nil, nil)
	errorhandler := e.ErrorHandler()
	if errorhandler == nil {
		t.Errorf("Error handler can't be nil")
	}
}
