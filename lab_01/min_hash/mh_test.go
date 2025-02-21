package min_hash

import (
	"testing"
)

func TestHashFunction_Hash(t *testing.T) {
	hf := NewHashFunction(42)
	hash := hf.Hash("test")
	if hash == 0 {
		t.Errorf("Expected non-zero hash, got %d", hash)
	}
}

func TestMinHash_Signature(t *testing.T) {
	mh := NewMinHash(0.1)
	set := []string{"apple", "banana", "cherry"}
	signature := mh.Signature(set)
	if len(signature) != mh.functionCount {
		t.Errorf("Expected signature length %d, got %d", mh.functionCount, len(signature))
	}
}

func TestMinHash_Similarity(t *testing.T) {
	mh := NewMinHash(0.1)
	setA := []string{"apple", "banana", "cherry"}
	setB := []string{"apple", "banana", "cherry"}
	sigA := mh.Signature(setA)
	sigB := mh.Signature(setB)
	similarity := mh.Similarity(sigA, sigB)
	if similarity != 1.0 {
		t.Errorf("Expected similarity 1.0, got %f", similarity)
	}
}

func TestMinHash_Similarity_DifferentSets(t *testing.T) {
	mh := NewMinHash(0.1)
	setA := []string{"apple", "banana", "cherry"}
	setB := []string{"apple", "banana", "date"}
	sigA := mh.Signature(setA)
	sigB := mh.Signature(setB)
	similarity := mh.Similarity(sigA, sigB)
	if similarity >= 1.0 {
		t.Errorf("Expected similarity less than 1.0, got %f", similarity)
	}
}

func BenchmarkHashFunction_Hash(b *testing.B) {
	hf := NewHashFunction(42)
	for i := 0; i < b.N; i++ {
		hf.Hash("benchmark")
	}
}

func BenchmarkMinHash_Signature(b *testing.B) {
	mh := NewMinHash(0.1)
	set := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mh.Signature(set)
	}
}

func BenchmarkMinHash_Similarity(b *testing.B) {
	mh := NewMinHash(0.1)
	setA := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew"}
	setB := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew"}
	sigA := mh.Signature(setA)
	sigB := mh.Signature(setB)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mh.Similarity(sigA, sigB)
	}
}

//todo дамп с википедией с историей правок страничек
