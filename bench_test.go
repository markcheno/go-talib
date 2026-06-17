package talib

import "testing"

// benchSeries builds a deterministic pseudo-random walk of length n for
// benchmarking window functions. Avoids math/rand so results are stable.
func benchSeries(n int) []float64 {
	s := make([]float64, n)
	v := 100.0
	for i := range s {
		// simple deterministic oscillation + drift
		v += float64((i*1103515245+12345)%97-48) / 50.0
		s[i] = v
	}
	return s
}

// MidPoint / MidPrice were rewritten from a naive O(n*period) full-window
// rescan to TA-Lib's amortized index-tracking approach. A large period makes
// the difference visible.

func BenchmarkMidPoint(b *testing.B) {
	in := benchSeries(5000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MidPoint(in, 200)
	}
}

func BenchmarkMidPrice(b *testing.B) {
	high := benchSeries(5000)
	low := make([]float64, len(high))
	for i := range high {
		low[i] = high[i] - 1.0
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MidPrice(high, low, 200)
	}
}
