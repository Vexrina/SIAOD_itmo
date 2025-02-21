package main

import (
	"fmt"
	"vexrina/siaod_itmo/lab_01/extendible_hashing"
	"vexrina/siaod_itmo/lab_01/min_hash"
)

func main() {
	mh_check()
}

func eh_check() {
	eh := extendible_hashing.NewExtendableHash(2, 5, true)
	eh.Insert("1", "1")
	eh.Insert("2", "1")
	eh.Insert("3", "1")
	eh.Insert("4", "1")
	eh.Insert("5", "1")
}

func mh_check() {
	setA := []string{"apple", "orange", "watermelon"}
	setB := []string{"apple", "orange", "pineapple"}

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
