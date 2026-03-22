package search

import (
	"encoding/binary"
	"math/rand"
	"strconv"
	"testing"

	ksearch "github.com/kelindar/search"
)

func setupTestIndex(size, dim int, seed int64) *ksearch.Index[[]byte] {
	rng := rand.New(rand.NewSource(seed))
	idx := ksearch.NewIndex[[]byte]()
	for i := 0; i < size; i++ {
		vec := make([]float32, dim)
		for j := range vec {
			vec[j] = rng.Float32()
		}
		lineID := i + 1
		payload := make([]byte, 4)
		binary.LittleEndian.PutUint32(payload, uint32(int32(lineID)))
		idx.Add(vec, payload)
	}
	return idx
}

func BenchmarkEngine_Search(b *testing.B) {
	idx := setupTestIndex(10000, 768, 42)
	engine := NewEngine(idx)
	rng := rand.New(rand.NewSource(99))

	queryVec := make([]float32, 768)
	for j := range queryVec {
		queryVec[j] = rng.Float32()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.Search(queryVec, 50, 10)
	}
}

func BenchmarkEngine_SearchTopKScaling(b *testing.B) {
	const (
		n   = 40000
		dim = 1000
	)

	idx := setupTestIndex(n, dim, 7)
	engine := NewEngine(idx)
	rng := rand.New(rand.NewSource(1001))

	queryVec := make([]float32, dim)
	for i := range queryVec {
		queryVec[i] = rng.Float32()
	}

	for _, k := range []int{25, 50, 100, 500, 1000, n} {
		b.Run("topK="+strconv.Itoa(k), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = engine.Search(queryVec, k, k/2)
			}
		})
	}
}
