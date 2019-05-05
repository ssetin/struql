package struql

import (
	"strings"
	"testing"
)

// Santa ...
type Santa struct {
	ID     int
	Clause string
}

// Origins ...
type Origins struct {
	Code  string
	Descr string
	Santa []Santa
}

// Device ...
type Device struct {
	Number      int
	Model       string
	Manufacture string
	Oi          []Origins
}

var (
	dev = Device{
		Number:      99,
		Model:       "JFQ",
		Manufacture: "Factory",
		Oi: []Origins{
			{
				Code:  "x256",
				Descr: "Debug",
				Santa: []Santa{
					{ID: 1, Clause: "Hoho"},
					{ID: 2, Clause: "Tree"},
					{ID: 3, Clause: "Tree"},
				},
			},
			{
				Code:  "x2599",
				Descr: "Release",
				Santa: []Santa{
					{ID: 1, Clause: "Tooo"},
					{ID: 2, Clause: "Mooo"},
					{ID: 3, Clause: "Laaa"},
				},
			},
		},
	}

	sq      StruQL
	filters = []Filter{
		{FieldName: "Oi.Descr", Value: "debug", Modifier: ModMe},
		{FieldName: "Oi.Santa.Clause", Value: "Tre", Operation: ComparsionBeginWith},
	}
)

func ModMe(s string) string {
	return strings.ToLower(s)
}

func init() {
	sq.Init(dev)
}

// BenchmarkIterateSearch ...
func BenchmarkIterateSearch(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for _, ois := range dev.Oi {
			if ModMe(ois.Descr) == "debug" {
				for _, snts := range ois.Santa {
					if strings.HasPrefix(snts.Clause, "Tre") {
						_ = snts.ID
					}
				}
			}
		}
	}
}

// BenchmarkSrtuqSearch ...
func BenchmarkSrtuqSearch(b *testing.B) {
	result := make(RowCollection, 0, 2)

	for i := 0; i < b.N; i++ {
		_, _ = sq.Where(result, filters...)
	}

}

// go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out
// go tool pprof -alloc_objects pkg.test mem.out
// go tool pprof pkg.test cpu.out

//BenchmarkIterateSearch-2         5000000               273 ns/op              32 B/op          4 allocs/op
//BenchmarkSrtuqSearch-2           1000000              1021 ns/op              80 B/op         12 allocs/op
