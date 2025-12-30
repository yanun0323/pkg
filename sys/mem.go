package sys

import "testing"

func MeasureMem(f func()) (allocs, bytes int64) {
	res := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f()
		}
	})
	return res.AllocsPerOp(), res.AllocedBytesPerOp()
}
