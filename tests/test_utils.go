package tests

import "math/rand/v2"

func generateRandomFloat32Array(n int) []float32 {
	arr := make([]float32, n)
	for i := range n {
		arr[i] = rand.Float32() // Generate random float64 between 0.0 and 1.0
	}
	return arr
}