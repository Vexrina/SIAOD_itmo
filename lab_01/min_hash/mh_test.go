package min_hash

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

var (
	oldDevilPath = "devil_05-12-2024"
	newDevilPath = "devil_21-03-2025"
	oldDevil     = readStringSlice(oldDevilPath)
	newDevil     = readStringSlice(newDevilPath)
)

func BenchmarkMinHash_Signature(b *testing.B) {
	mh := NewMinHash(0.1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mh.Signature(oldDevil)
	}
}

func BenchmarkMinHash_Similarity(b *testing.B) {
	mh := NewMinHash(0.1)
	sigA := mh.Signature(oldDevil)
	sigB := mh.Signature(newDevil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mh.Similarity(sigA, sigB)
	}
}

func readStringSlice(filePath string) []string {
	//dir, err := os.Getwd()
	//if err != nil {
	//	panic(fmt.Sprint("Ошибка получения рабочей директории:", err))
	//}
	//fmt.Println("Текущая директория:", dir)
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprint("Ошибка открытия файла:", err))
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var words []string
	for scanner.Scan() {
		lineWords := strings.Fields(scanner.Text())
		words = append(words, lineWords...)
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprint("Ошибка чтения файла:", err))
		return nil
	}

	return words
}
