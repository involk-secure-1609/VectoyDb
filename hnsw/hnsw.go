package hnsw

import (
	"cmp"
	"math"
	"math/rand"
	"time"
)

type Node[K cmp.Ordered] struct {
	key        K
	embedding  []float64
	neighbours []Node[K]
}
type level[K cmp.Ordered] struct {
	nodes map[K]Node[K]
}

// EuclideanDistSquare calculates the squared Euclidean distance between two vectors.
// It's computationally cheaper than Euclidean distance and often sufficient for comparisons.
func EuclideanDistSquare(p1 []float64, p2 []float64) float64 {
	var sum float64 = 0
	for i := range p1 {
		d := p2[i] - p1[i] // Calculate the difference between corresponding coordinates.
		sum += d * d       // Square the difference and add to the sum.
	}
	return sum
}

// EuclideanDistance computes the Euclidean distance between two vectors.
func EuclideanDistance(vec1 []float64, vec2 []float64) float64 {
	// TODO: can we speedup with vek?
	var sum float64 = 0
	for i := range vec1 {
		diff := vec1[i] - vec2[i]
		sum += diff * diff
	}
	return float64(math.Sqrt(float64(sum)))
}

// DotProduct computes the DotProduct distance between two vectors.
func DotProduct(vec1 []float64, vec2 []float64) float64 {
	// TODO: can we speedup with vek?
	var sum float64 = 0
	for i := range vec1 {
		prod := vec1[i] * vec2[i]
		sum += prod
	}
	return sum
}

var distanceFuncs = map[string]DistanceFunc{
	"euclidean":  EuclideanDistance,
	"dotProduct": DotProduct,
	"squareDistance": EuclideanDistSquare,
}

type DistanceFunc func(vec1 []float64, vec2 []float64) float64

type HNSWGraph[K cmp.Ordered] struct {
	// Distance is the distance function used to compare embeddings.
	Distance DistanceFunc
	Rng      *rand.Rand
	// M is the maximum number of neighbors to keep for each node.
	// A good default for OpenAI embeddings is 16.
	M int

	// Ml is the level generation factor.
	// E.g., for Ml = 0.25, each layer is 1/4 the size of the previous layer.
	Ml float64

	// EfSearch is the number of nodes to consider in the search phase.
	// 20 is a reasonable default. Higher values improve search accuracy at
	// the expense of memory.
	EfSearch int

	levels []level[K]
}

func NewHNSWGraph[K cmp.Ordered](distanceFunc string) *HNSWGraph[K] {

	return &HNSWGraph[K]{
		M:        16,
		Ml:       0.25,
		EfSearch: 20,
		Rng:      defaultRand(),
		Distance: distanceFuncs[distanceFunc],
	}
}

func defaultRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
// maxLevel returns an upper-bound on the number of levels in the graph
// based on the size of the base layer.
func maxLevel(ml float64, numNodes int) int {
	if ml == 0 {
		panic("ml must be greater than 0")
	}

	if numNodes == 0 {
		return 1
	}

	l := math.Log(float64(numNodes))
	l /= math.Log(1 / ml)

	m := int(math.Round(l)) + 1

	return m
}

// randomLevel generates a random level for a new node.
func (g *HNSWGraph[K]) randomLevel() int {
	// max avoids having to accept an additional parameter for the maximum level
	// by calculating a probably good one from the size of the base layer.
	max := 1
	if len(g.levels) > 0 {
		if g.Ml == 0 {
			panic("(*Graph).Ml must be greater than 0")
		}
		//max = maxLevel(g.Ml, g.levels[0])
	}

	// for level := 0; level < max; level++ {
	// 	if h.Rng == nil {
	// 		h.Rng = defaultRand()
	// 	}
	// 	r := h.Rng.Float64()
	// 	if r > h.Ml {
	// 		return level
	// 	}
	// }
	
	return max
}
func (g *HNSWGraph[K]) insert(key K,embedding []float64){

}