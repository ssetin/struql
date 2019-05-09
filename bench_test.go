package struql

import (
	"strings"
	"testing"
)

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
	filters := []Filter{
		{FieldName: "Oi.Descr", Value: "Debug"},
		{FieldName: "Oi.Santa.Clause", Value: "Tre", Operation: ComparisonBeginWith},
	}
	result := make(RowCollection, 0, 2)

	for i := 0; i < b.N; i++ {
		_, _ = sq.Where(result, filters...)
	}
}

func BenchmarkIterateSearchWithModifier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, ois := range dev.Oi {
			if strings.ToLower(ois.Descr) == "debug" {
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
func BenchmarkSrtuqSearchWithModifier(b *testing.B) {
	filters := []Filter{
		{FieldName: "Oi.Descr", Value: "debug", Modifier: modMe},
		{FieldName: "Oi.Santa.Clause", Value: "Tre", Operation: ComparisonBeginWith},
	}
	result := make(RowCollection, 0, 2)

	for i := 0; i < b.N; i++ {
		_, _ = sq.Where(result, filters...)
	}

}

// go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out
// go tool pprof -alloc_objects struql.test mem.out
// go tool pprof struql.test cpu.out

//BenchmarkIterateSearch-2                20000000              76.5 ns/op               0 B/op          0 allocs/op
//BenchmarkSrtuqSearch-2                   1000000              1080 ns/op               0 B/op          0 allocs/op
//BenchmarkIterateSearchWithModifier-2     5000000               366 ns/op              48 B/op          6 allocs/op
//BenchmarkSrtuqSearchWithModifier-2        300000              3808 ns/op             560 B/op         54 allocs/op
