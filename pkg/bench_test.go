package struql

import (
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
	Number int
	Model  string
	Man    string
	Oi     Origins
}

var (
	dev = Device{
		Number: 99,
		Model:  "JFQ",
		Man:    "Sony",
		Oi: Origins{
			Code:  "x256",
			Descr: "Debug",
			Santa: []Santa{
				{ID: 1, Clause: "Hoho"},
				{ID: 2, Clause: "Rum"},
				{ID: 3, Clause: "Tree"},
			},
		}}

	sq     StruQL
	filter = Filter{FieldName: "Oi.Santa.Clause", Value: "Rum"}
)

func init() {
	sq.Init(dev)
}

// BenchmarkIterateSearch ...
func BenchmarkIterateSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for idx := 0; idx < len(dev.Oi.Santa); idx++ {
			if dev.Oi.Santa[idx].Clause == "Rum" {
				_ = dev.Oi.Santa[idx]
			}
		}
	}
}

// BenchmarkSrtuqSearch ...
func BenchmarkSrtuqSearch(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_ = sq.Where(filter)
	}
}

// go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out
// go tool pprof pkg.test mem.out

//BenchmarkIterateSearch-2        100000000               10.0 ns/op             0 B/op          0 allocs/op
//BenchmarkSrtuqSearch-2            2000000                743 ns/op           152 B/op          4 allocs/op
//BenchmarkSrtuqSearch-2            5000000                269 ns/op            24 B/op          2 allocs/op
//BenchmarkSrtuqSearch-2           10000000                188 ns/op            16 B/op          1 allocs/op
