package tests

import "math/rand/v2"

func generateRandomFloat64Array(n int) []float64 {
	arr := make([]float64, n)
	for i := range n {
		arr[i] = rand.Float64() // Generate random float64 between 0.0 and 1.0
	}
	return arr
}