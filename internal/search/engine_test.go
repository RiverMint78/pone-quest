package search

import (
	"math/rand"
	"testing"

	ksearch "github.com/kelindar/search"
)

func setupTestIndex(size int) *ksearch.Index[string] {
	idx := ksearch.NewIndex[string]()
	for i := 0; i < size; i++ {
		vec := make([]float32, 768)
		for j := range vec {
			vec[j] = rand.Float32()
		}
		idx.Add(vec, "test_id")
	}
	return idx
}

func BenchmarkEngine_Search(b *testing.B) {
	idx := setupTestIndex(10000)
	engine := NewEngine(idx)

	queryVec := make([]float32, 768)
	for j := range queryVec {
		queryVec[j] = rand.Float32()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.Search(queryVec, 50)
	}
}
