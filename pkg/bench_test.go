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

	sq     StruQL
	filter = []Filter{
		{FieldName: "Oi.Descr", Value: "Debug"},
		{FieldName: "Oi.Santa.Clause", Value: "Tre", Operation: ComparsionBeginWith},
	}
)

func init() {
	sq.Init(dev)
}

// BenchmarkIterateSearch ...
func BenchmarkIterateSearch(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for _, ois := range dev.Oi {
			if ois.Descr == "Debug" {
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
		_, _ = sq.Where(result, filter...)
	}

}

// go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out
// go tool pprof -alloc_objects pkg.test mem.out
// go tool pprof pkg.test cpu.out

//BenchmarkIterateSearch-2        100000000               10.0 ns/op             0 B/op          0 allocs/op
//BenchmarkSrtuqSearch-2            2000000                596 ns/op            16 B/op          1 allocs/op
//BenchmarkSrtuqSearch-2            3000000                607 ns/op            16 B/op          1 allocs/op
//BenchmarkSrtuqSearch-2            3000000                440 ns/op             0 B/op          0 allocs/op

//BenchmarkIterateSearch-4        100000000               15.4 ns/op             0 B/op          0 allocs/op
//BenchmarkSrtuqSearch-4            5000000                237 ns/op             0 B/op          0 allocs/op
