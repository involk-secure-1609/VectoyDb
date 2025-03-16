package hnsw

import (
	"cmp"
	"fmt"
	"maps"
	"math"
	"math/rand"
	"slices"
	"time"
)

type Embedding []float64
type Node[K cmp.Ordered] struct {
	Key        K
	Embed      Embedding
	neighbours map[K]*Node[K]
}

func MakeNode[K cmp.Ordered](key K, embed Embedding) Node[K] {
	return Node[K]{Key: key, Embed: embed}
}

type DistanceFunc func(vec1 Embedding, vec2 Embedding) float64

// EuclideanDistance computes the Euclidean distance between two vectors.
func EuclideanDistance(a, b Embedding) float64 {
	// TODO: can we speedup with vek?
	var sum float64 = 0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return float64(math.Sqrt(float64(sum)))
}

// EuclideanDistSquare calculates the squared Euclidean distance between two vectors.
// It's computationally cheaper than Euclidean distance and often sufficient for comparisons.
func EuclideanDistSquare(p1 Embedding, p2 Embedding) float64 {
	var sum float64 = 0
	for i := range p1 {
		d := p2[i] - p1[i] // Calculate the difference between corresponding coordinates.
		sum += d * d       // Square the difference and add to the sum.
	}
	return sum
}

// DotProduct computes the DotProduct distance between two vectors.
func DotProduct(vec1 Embedding, vec2 Embedding) float64 {
	// TODO: can we speedup with vek?
	var sum float64 = 0
	for i := range vec1 {
		prod := vec1[i] * vec2[i]
		sum += prod
	}
	return sum
}

var distanceFuncs = map[string]DistanceFunc{
	"euclidean":      EuclideanDistance,
	"dotProduct":     DotProduct,
	"squareDistance": EuclideanDistSquare,
}


type searchCandidate[K cmp.Ordered] struct {
	node *Node[K]
	dist float64
}

func (s searchCandidate[K]) Less(o searchCandidate[K]) bool {
	return s.dist < o.dist
}

// search returns the node closest to the target node
// within the same level.
func (n *Node[K]) search(
	// k is the number of candidates in the result set.
	k int,
	efSearch int,
	target Embedding,
	distance DistanceFunc,
) []searchCandidate[K] {
	// This is a basic greedy algorithm to find the entry point at the given level
	// that is closest to the target node.
	candidates := Heap[searchCandidate[K]]{}
	candidates.Init(make([]searchCandidate[K], 0, efSearch))
	candidates.Push(
		searchCandidate[K]{
			node: n,
			dist: distance(n.Embed, target),
		},
	)
	var (
		result  = Heap[searchCandidate[K]]{}
		visited = make(map[K]bool)
	)
	result.Init(make([]searchCandidate[K], 0, k))

	// Begin with the entry node in the result set.
	result.Push(candidates.Min())
	visited[n.Key] = true

	for candidates.Len() > 0 {
		var (
			current  = candidates.Pop().node
			improved = false
		)

		// We iterate the map in a sorted, deterministic fashion for
		// tests.
		neighborKeys := maps.Keys(current.neighbours)
		sortedNeighborKeys:=slices.Sorted(neighborKeys)
		for _, neighborID := range sortedNeighborKeys {
			neighbor := current.neighbours[neighborID]
			if visited[neighborID] {
				continue
			}
			visited[neighborID] = true

			dist := distance(neighbor.Embed, target)
			improved = improved || dist < result.Min().dist
			if result.Len() < k {
				result.Push(searchCandidate[K]{node: neighbor, dist: dist})
			} else if dist < result.Max().dist {
				result.PopLast()
				result.Push(searchCandidate[K]{node: neighbor, dist: dist})
			}

			candidates.Push(searchCandidate[K]{node: neighbor, dist: dist})
			// Always store candidates if we haven't reached the limit.
			if candidates.Len() > efSearch {
				candidates.PopLast()
			}
		}

		// Termination condition: no improvement in distance and at least
		// kMin candidates in the result set.
		if !improved && result.Len() >= k {
			break
		}
	}

	return result.Slice()
}


func (node *Node[K]) addNeighbour(neighbor *Node[K], m int, dist DistanceFunc) {
	if node.neighbours == nil {
		node.neighbours = make(map[K]*Node[K])
	}

	node.neighbours[neighbor.Key] = neighbor
	if len(node.neighbours) <= m {
		return
	}

	// Find the neighbor with the worst distance.
	var (
		worstDist = float64(math.Inf(-1))
		worst     *Node[K]
	)
	for _, neighbor := range node.neighbours {
		d := dist(neighbor.Embed, node.Embed)
		// d > worstDist may always be false if the distance function
		// returns NaN, e.g., when the embeddings are zero.
		if d > worstDist || worst == nil {
			worstDist = d
			worst = neighbor
		}
	}

	delete(node.neighbours, worst.Key)
	// Delete backlink from the worst neighbor.
	delete(worst.neighbours, node.Key)
	worst.replenish(m)
}

func (node *Node[K]) replenish(m int) {
	if len(node.neighbours) >= m {
		return
	}
	// Restore connectivity by adding new neighbors.
	// This is a naive implementation that could be improved by
	// using a priority queue to find the best candidates.
	for _, neighbour := range node.neighbours {
		for key, candidate := range neighbour.neighbours {
			if _, ok := node.neighbours[key]; ok {
				// do not add duplicates
				continue
			}
			if candidate == node {
				continue
			}
			node.addNeighbour(candidate, m, EuclideanDistSquare)
			if len(node.neighbours) >= m {
				return
			}
		}
	}
}

// isolates remove the node from the graph by removing all connections
// to neighbors.
func (node *Node[K]) isolate(m int) {
	for _, neighbor := range node.neighbours {
		delete(neighbor.neighbours, node.Key)
		neighbor.replenish(m)
	}
}

type level[K cmp.Ordered] struct {
	// nodes is a map of nodes IDs to nodes.
	// All nodes in a higher level are also found in the lower levels,
	// an essential property of the graph.
	//
	// nodes is exported for interop with encoding/gob.
	nodes map[K]*Node[K]
}

// entry returns the entry node of the layer.
// It doesn't matter which node is returned, even that the
// entry node is consistent, so we just return the first node
// in the map to avoid tracking extra state.
func (l *level[K]) entry() *Node[K] {
	if l == nil {
		return nil
	}
	for _, node := range l.nodes {
		return node
	}
	return nil
}

func (l *level[K]) size() int {
	if l == nil {
		return 0
	}
	return len(l.nodes)
}

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

	levels []*level[K]
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
		max = maxLevel(g.Ml, g.levels[0].size())
	}

	for level := range max {
		if g.Rng == nil {
			g.Rng = defaultRand()
		}
		r := g.Rng.Float64()
		if r > g.Ml {
			return level
		}
	}

	return max
}

func (g *HNSWGraph[K]) assertDims(n Embedding) {
	if len(g.levels) == 0 {
		return
	}
	hasDims := g.Dims()
	if hasDims != len(n) {
		panic(fmt.Sprint("embedding dimension mismatch: ", hasDims, " != ", len(n)))
	}
}

// Dims returns the number of dimensions in the graph, or
// 0 if the graph is empty.
func (g *HNSWGraph[K]) Dims() int {
	if len(g.levels) == 0 {
		return 0
	}
	return len(g.levels[0].entry().Embed)
}

// Len returns the number of nodes in the graph.
func (g *HNSWGraph[K]) Len() int {
	if len(g.levels) == 0 {
		return 0
	}
	return g.levels[0].size()
}

func ptr[T any](v T) *T {
	return &v
}

// inserts nodes into the graph.
// If another node with the same ID exists, it is replaced.
func (g *HNSWGraph[K]) Insert(nodes ...Node[K]) {
	for _, node := range nodes {
		key := node.Key
		embedding := node.Embed

		g.assertDims(embedding)
		insertLevel := g.randomLevel()
		// Create layers that don't exist yet.
		for insertLevel >= len(g.levels) {
			g.levels = append(g.levels, &level[K]{})
		}

		if insertLevel < 0 {
			panic("invalid level")
		}

		var elevator *K

		preLen := g.Len()

		// Insert node at each level, beginning with the highest.
		for i := len(g.levels) - 1; i >= 0; i-- {
			level := g.levels[i]
			newNode :=&Node[K]{
				Key: key,
				Embed: embedding,
			}
			

			// Insert the new node into the layer.
			if level.entry() == nil {
				level.nodes = make(map[K]*Node[K])
				level.nodes[key]=newNode
				continue
			}

			// Now at the highest level with more than one node, so we can begin
			// searching for the best way to enter the graph.
			searchPoint := level.entry()

			// On subsequent layers, we use the elevator node to enter the graph
			// at the best point.
			if elevator != nil {
				searchPoint = level.nodes[*elevator]
			}

			if g.Distance == nil {
				panic("(*Graph).Distance must be set")
			}

			neighborhood := searchPoint.search(g.M, g.EfSearch, embedding, g.Distance)
			if len(neighborhood) == 0 {
				// This should never happen because the searchPoint itself
				// should be in the result set.
				panic("no nodes found")
			}

			// Re-set the elevator node for the next layer.
			elevator = ptr(neighborhood[0].node.Key)

			if insertLevel >= i {
				if _, ok := level.nodes[key]; ok {
					g.Delete(key)
				}
				// Insert the new node into the layer.
				level.nodes[key] = newNode
				for _, node := range neighborhood {
					// Create a bi-directional edge between the new node and the best node.
					node.node.addNeighbour(newNode, g.M, g.Distance)
					newNode.addNeighbour(node.node, g.M, g.Distance)
				}
			}
		}

		// Invariant check: the node should have been added to the graph.
		if g.Len() != preLen+1 {
			panic("node not added")
		}
	}
}


// Search finds the k nearest neighbors from the target node.
func (h *HNSWGraph[K]) Search(near Embedding, k int) []Node[K] {
	h.assertDims(near)
	if len(h.levels) == 0 {
		return nil
	}

	var (
		efSearch = h.EfSearch

		elevator *K
	)

	for level := len(h.levels) - 1; level >= 0; level-- {
		searchPoint := h.levels[level].entry()
		if elevator != nil {
			searchPoint = h.levels[level].nodes[*elevator]
		}

		// Descending hierarchies
		if level > 0 {
			nodes := searchPoint.search(1, efSearch, near, h.Distance)
			elevator = ptr(nodes[0].node.Key)
			continue
		}

		nodes := searchPoint.search(k, efSearch, near, h.Distance)
		out := make([]Node[K], 0, len(nodes))

		for _, node := range nodes {
			out = append(out, *node.node)
		}

		return out
	}

	panic("unreachable")
}

// Delete removes a node from the graph by key.
// It tries to preserve the clustering properties of the graph by
// replenishing connectivity in the affected neighborhoods.
func (h *HNSWGraph[K]) Delete(key K) bool {
	if len(h.levels) == 0 {
		return false
	}

	var deleted bool
	for _, layer := range h.levels {
		node, ok := layer.nodes[key]
		if !ok {
			continue
		}
		delete(layer.nodes, key)
		node.isolate(h.M)
		deleted = true
	}

	return deleted
}

// Lookup returns the vector with the given key.
func (h *HNSWGraph[K]) Lookup(key K) (Embedding, bool) {
	if len(h.levels) == 0 {
		return nil, false
	}

	node, ok := h.levels[0].nodes[key]
	if !ok {
		return nil, false
	}
	return node.Embed, ok
}