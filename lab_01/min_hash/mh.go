package min_hash

import (
	"math/rand"
)

// HashFunction представляет хэш-функцию
type HashFunction struct {
	seed int
}

// NewHashFunction создает новую хэш-функцию с заданным seed
func NewHashFunction(seed int) *HashFunction {
	return &HashFunction{seed: seed}
}

// Hash вычисляет хэш для строки
func (hf *HashFunction) Hash(s string) uint32 {
	result := uint32(1)
	for _, char := range s {
		result = uint32(hf.seed)*result + uint32(char)
	}
	return result
}

// MinHash представляет структуру для MinHash
type MinHash struct {
	functionCount int
	functions     []*HashFunction
}

// NewMinHash создает новый MinHash с заданной максимальной ошибкой
func NewMinHash(maxError float64) *MinHash {
	functionCount := int(1 / (maxError * maxError))
	functions := make([]*HashFunction, functionCount)

	for i := 0; i < functionCount; i++ {
		seed := rand.Intn(functionCount) + 32
		functions[i] = NewHashFunction(seed)
	}

	return &MinHash{
		functionCount: functionCount,
		functions:     functions,
	}
}

// FindMin находит минимальный хэш для множества с использованием заданной хэш-функции
func (mh *MinHash) FindMin(set []string, hf *HashFunction) uint32 {
	minHash := uint32(0xFFFFFFFF)
	for _, s := range set {
		hash := hf.Hash(s)
		if hash < minHash {
			minHash = hash
		}
	}
	return minHash
}

// Signature вычисляет сигнатуру для множества
func (mh *MinHash) Signature(set []string) []uint32 {
	signature := make([]uint32, mh.functionCount)
	for i, hf := range mh.functions {
		signature[i] = mh.FindMin(set, hf)
	}
	return signature
}

// Similarity вычисляет схожесть между двумя сигнатурами
func (mh *MinHash) Similarity(sigA, sigB []uint32) float64 {
	equalCount := 0
	for i := 0; i < mh.functionCount; i++ {
		if sigA[i] == sigB[i] {
			equalCount++
		}
	}
	return float64(equalCount) / float64(mh.functionCount)
}
