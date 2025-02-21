package main

import (
	"fmt"
	"vexrina/siaod_itmo/lab_01/min_hash"
)

func main() {
	// Пример использования
	setA := []string{"apple", "orange"}
	setB := []string{"apple", "peach"}

	sumSim := 0.0
	counter := 0
	for i := 0; i < 10000; i++ {
		minHash := min_hash.NewMinHash(0.1)
		sigA := minHash.Signature(setA)
		sigB := minHash.Signature(setB)
		counter++
		sumSim += minHash.Similarity(sigA, sigB)
	}

	fmt.Printf("Similarity: %.2f\n", sumSim/float64(counter))
}
