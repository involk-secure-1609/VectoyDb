package lsh

import (
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

// hyperplanes represents a collection of hyperplanes.
// Each hyperplane is a vector in the same dimensional space as the input data points.
type hyperplanes [][]float64

// Key is a way to index into a table. In this context, it's a binary hash key
// represented as a slice of uint8 (0 or 1).
type hashTableKey []uint8

// Point represents an abstract point in n-dimensional space.
// It contains the vector itself, any extra data associated with the point, and a unique ID.
type Point struct {
	Vector    []float64 // The vector representing the point in n-dimensional space.
	ExtraData string    // Optional extra data associated with the point.
	ID        uint64    // Unique identifier for the point.
}

// QueryResult represent query result with distance to query point.
// It includes the Point itself and the calculated distance to the query point.
type QueryResult struct {
	Point            // The Point that is a query result.
	Distance float64 // The distance between the query point and this Point.
}

// NewHyperplanes generates and initializes a set of d hyperplanes with s dimensions.
// d is the number of hyperplanes to generate.
// s is the number of dimensions each hyperplane will have (same as input data points).
func newHyperplanes(d, s int32) hyperplanes {
	hs := make([][]float64, d) // Create a slice of slices to hold 'd' hyperplanes.
	for i := range d {
		v := make([]float64, s) // Each hyperplane is a vector of 's' dimensions.
		for j := range s {
			n := rand.NormFloat64() // Generate a random number from a normal (Gaussian) distribution.
			v[j] = n                // Assign this random number as a coordinate in the hyperplane vector.
		}
		hs[i] = v // Add the generated hyperplane vector to the set of hyperplanes.
	}
	return hs
}

// DistanceFunc is a function type for calculating the distance between two vectors.
// It takes two float64 slices (representing vectors) and returns a float64 (the distance).
type DistanceFunc func(p1 []float64, p2 []float64) float64

// euclideanDistSquare calculates the squared Euclidean distance between two vectors.
// It's computationally cheaper than Euclidean distance and often sufficient for comparisons.
func euclideanDistSquare(p1 []float64, p2 []float64) (sum float64) {
	for i := range p1 {
		d := p2[i] - p1[i] // Calculate the difference between corresponding coordinates.
		sum += d * d       // Square the difference and add to the sum.
	}
	return sum
}

var dFuncMap = map[string]DistanceFunc{
	"euclidean": euclideanDistSquare,
}

// hashTableBucket is a bucket in the hash table.
// It's a slice of Points that hash to the same key.
type hashTableBucket []Point

// hashTable is the hash table data structure.
// It's a map where the key is a uint64 (the hash key) and the value is a hashTableBucket.
type hashTable map[uint64]hashTableBucket

// signature represents a simhash signature - an array of hash values (0s and 1s).
type signature []uint8

// simhash struct holds the signature.
type simhash struct {
	sig signature // The signature (array of 0s and 1s).
}

// newSimhash generates the simhash of an attribute (a vector) using the hyperplanes.
// hs: The set of hyperplanes to use.
// e: The input vector for which to generate the simhash.
func newSimhash(hs hyperplanes, e []float64) *simhash {
	sig := newSignature(hs, e) // Generate the signature using the hyperplanes and input vector.
	return &simhash{
		sig: sig, // Create and return a simhash struct with the generated signature.
	}
}

// newSignature computes the signature for a simhash of input float array.
// hyperplanes: The set of hyperplanes.
// e: The input vector.
func newSignature(hyperplanes hyperplanes, e []float64) signature {
	sigarr := make([]uint8, len(hyperplanes)) // Initialize a slice to store the signature bits (0 or 1).
	for hix, h := range hyperplanes {         // Iterate through each hyperplane.
		var dp float64        // Initialize a variable to store the dot product.
		for k, v := range e { // Calculate the dot product of the hyperplane and the input vector.
			dp += h[k] * float64(v) // Multiply corresponding coordinates and sum them up.
		}
		if dp >= 0 { // If the dot product is non-negative, the point is on one side of the hyperplane.
			sigarr[hix] = uint8(1) // Assign 1 to the signature bit.
		} else { // Otherwise, the point is on the other side.
			sigarr[hix] = uint8(0) // Assign 0 to the signature bit.
		}
	}
	return sigarr // Return the generated signature.
}

// Hash returns all combined hash values for all hash tables.
// For each hash table, it generates a hash key based on the simhash and LSH parameters.
func (clsh *cosineLshParam) hash(point []float64) []hashTableKey {
	simhash := newSimhash(clsh.hyperplanes, point) // Generate the simhash for the input point.
	hvs := make([]hashTableKey, clsh.l)            // Create a slice to hold hash keys for each of the 'l' hash tables.
	for i := range hvs {                           // For each hash table...
		s := make(hashTableKey, clsh.m) // Create a hash key of length 'm' (number of hash functions per table).
		for j := range clsh.m {         // For each hash function in this hash table...
			s[j] = uint8(simhash.sig[int32(i)*clsh.m+j]) // Take 'm' bits from the simhash signature to form the hash key.
			// i*clsh.m+j:  Calculates the index into the simhash.sig array.
			// For table 0, it takes bits 0 to m-1.
			// For table 1, it takes bits m to 2m-1, and so on.
		}
		hvs[i] = s // Assign the generated hash key to the slice of hash keys for this hash table.
	}
	return hvs // Return the slice of hash keys, one for each hash table.
}

// CosineLsh is an implementation of Random projection LSH (Locality-Sensitive Hashing).
// https://en.wikipedia.org/wiki/Locality-sensitive_hashing#Random_projection
type CosineLsh struct {
	*cosineLshParam             // Embed the LSH parameters.
	tables          []hashTable // Slice of hash tables, each is a map.
	nextID          uint64      // Atomic counter to generate unique IDs for inserted points.
}

// cosineLshParam holds the parameters for Cosine LSH.
type cosineLshParam struct {
	dim         int32        // Dimensionality of the input data points.
	l           int32         // Number of hash tables to use.
	m           int32         // Number of hash functions (hyperplanes) used in each hash table to create a hash key.
	hyperplanes [][]float64 // The set of randomly generated hyperplanes used for hashing.
	h           int32         // Total number of hyperplanes (l * m).
	dFunc       string      // Function to calculate the distance between vectors.
}

// NewLshParams initializes the LSH settings.
// dim: Dimensionality of input vectors.
// l: Number of hash tables.
// m: Number of hash functions per table.
// h: Total number of hash functions.
// hyperplanes: Pre-generated hyperplanes.
func newCosineLshParam(dim, l, m, h int32, dFunc string, hyperplanes [][]float64) *cosineLshParam {
	return &cosineLshParam{
		dim:         dim,         // Set the dimensionality.
		l:           l,           // Set the number of hash tables.
		m:           m,           // Set the number of hash functions per table.
		hyperplanes: hyperplanes, // Assign the generated hyperplanes.
		h:           h,           // Set the total number of hyperplanes.
		dFunc:       dFunc,       // Use squared Euclidean distance as the default distance function.
	}
}


// NewCosineLsh creates an instance of Cosine LSH.
// dim is the number of dimensions of the input points.
// l is the number of hash tables.
// m is the number of hash values in each hash table (length of hash key for each table).
func NewCosineLsh(dim, l, m int32, dfunc string) *CosineLsh {
	h := m * l                            // Calculate the total number of hyperplanes needed.
	hyperplanes := newHyperplanes(h, dim) // Generate 'h' hyperplanes of 'dim' dimensions.
	tables := make([]hashTable, l)        // Create 'l' hash tables.
	for i := range tables {
		tables[i] = make(hashTable) // Initialize each hash table as an empty map.
	}
	return &CosineLsh{
		cosineLshParam: newCosineLshParam(dim, l, m, h, dfunc, hyperplanes), // Initialize LSH parameters.
		tables:         tables,                                              // Assign the created hash tables.
	}
}

// Insert adds a new data point to the Cosine LSH index.
// point is the data point (vector) to be inserted.
// extraData is any additional data to be stored with the point.
func (lsh *CosineLsh) Insert(point []float64, extraData string) {
	// Apply hash functions to generate hash keys for the point in each hash table.
	hvs := lsh.toBasicHashTableKeys(lsh.hash(point))
	// Insert the point into all hash tables.
	var wg sync.WaitGroup         // WaitGroup to manage concurrent insertions into hash tables.
	wg.Add(len(lsh.tables))     // Increment WaitGroup counter by the number of hash tables.
	for i := range lsh.tables { // Iterate through each hash table.
		hv := hvs[i]                          // Get the hash key for the current hash table.
		table := lsh.tables[i]              // Get the current hash table.
		go func(table hashTable, hv uint64) { // Launch a goroutine for concurrent insertion.
			defer wg.Done()                    // Decrement WaitGroup counter when the goroutine finishes.
			if _, exist := table[hv]; !exist { // Check if a bucket exists for this hash key in the table.
				table[hv] = make(hashTableBucket, 0) // If not, create a new empty bucket.
			}
			vectorID := atomic.AddUint64(&lsh.nextID, 1)                                          // Atomically increment the point ID counter and get the new ID.
			table[hv] = append(table[hv], Point{Vector: point, ID: vectorID, ExtraData: extraData}) // Append the point to the bucket associated with the hash key.
		}(table, hv)
	}
	wg.Wait() // Wait for all goroutines to complete before returning.
}

// Delete removes a new data point from the Cosine LSH index.
// point is the data point (vector) to be removed.
// extraData is any additional data which is stored with the point.
func (lsh *CosineLsh) Delete(point []float64, extraData string) {
	// Apply hash functions to generate hash keys for the point in each hash table.
	hvs := lsh.toBasicHashTableKeys(lsh.hash(point))
	var wg sync.WaitGroup         // WaitGroup to manage concurrent insertions into hash tables.
	wg.Add(len(lsh.tables))     // Increment WaitGroup counter by the number of hash tables.
	for i := range lsh.tables { // Iterate through each hash table.
		hv := hvs[i]                          // Get the hash key for the current hash table.
		table := lsh.tables[i]              // Get the current hash table.
		go func(table hashTable, hv uint64, point []float64, extraData string) { // Launch a goroutine for concurrent deletion.
			defer wg.Done() // Decrement WaitGroup counter when the goroutine finishes.
		
			// Check if a bucket exists for this hash key in the table.
			if _, exist := table[hv]; !exist {
				return
			}
			
			// Find and remove the matching point
			var newBucket []Point
			for _, p := range table[hv] {
				// Skip points that match both vector and extraData (effectively deleting them)
				if vectorsEqual(p.Vector, point) && p.ExtraData == extraData {
					continue
				}
				// Keep all other points
				newBucket = append(newBucket, p)
			}
			
			// Replace the bucket with the filtered version
			table[hv] = newBucket
		}(table, hv, point, extraData)
	}
	wg.Wait() // Wait for all goroutines to complete before returning.
}

// Helper function to check if two vectors are equal
func vectorsEqual(a, b []float64) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

// Query finds the approximate nearest neighbors of a query point.
// q is the query point (vector).
// maxResult is the maximum number of results to return (if > 0, returns top 'maxResult' nearest neighbours).
func (lsh *CosineLsh) Search(q []float64, maxResult int) []QueryResult {
	// Apply hash functions to the query point to get hash keys for each hash table.
	hvs := lsh.toBasicHashTableKeys(lsh.hash(q))
	// Keep track of points seen to avoid duplicates (across different hash tables).
	seen := make(map[uint64]Point)       // Map to store unique points, keyed by their IDs.
	for i, table := range lsh.tables { // Iterate through each hash table and corresponding hash key.
		if candidates, exist := table[hvs[i]]; exist { // Check if a bucket exists in the current table for the hash key.
			for _, id := range candidates { // Iterate through the points in the bucket (candidates).
				if _, exist := seen[id.ID]; exist { // Check if this point has already been seen (processed from another table).
					continue // If seen, skip to the next point to avoid duplicates.
				}
				seen[id.ID] = id // If not seen, add the point to the 'seen' map.
			}
		}
	}

	distances := make([]QueryResult, 0, len(seen)) // Create a slice to store QueryResults.
	for _, value := range seen {                   // Iterate through the unique points found in the hash tables.
		dist := dFuncMap[lsh.dFunc](q, value.Vector) // Calculate the distance between the query point and the candidate point.
		queryResult := QueryResult{Distance: dist}     // Create a QueryResult struct.
		queryResult.Point = value                      // Assign the Point to the QueryResult.
		distances = append(distances, queryResult)     // Add the QueryResult to the slice.
	}
	sort.Slice(distances, func(i, j int) bool { // Sort the results by distance in ascending order (nearest first).
		return distances[i].Distance < distances[j].Distance
	})

	if maxResult > 0 && len(distances) > maxResult { // If maxResult is specified and there are more results than maxResult.
		return distances[:maxResult] // Return only the top 'maxResult' nearest neighbors.
	}

	return distances // Return all QueryResults if maxResult is not specified or if there are fewer results.
}

// toBasicHashTableKeys converts hashTableKey (slice of uint8) to uint64.
// This is needed because Go maps use comparable keys, and slices are not directly comparable.
// Converting the binary hash key (slice of 0s and 1s) to a uint64 allows it to be used as a map key.
func (lsh *CosineLsh) toBasicHashTableKeys(keys []hashTableKey) []uint64 {
	basicKeys := make([]uint64, lsh.cosineLshParam.l) // Create a slice to store the converted uint64 keys.
	for i, key := range keys {                          // Iterate through the slice of hashTableKeys.
		s := ""                       // Initialize an empty string to build the binary string representation.
		for _, hashVal := range key { // Iterate through the hash key (slice of uint8).
			switch hashVal {
			case uint8(0): // If the hash bit is 0.
				s += "0" // Append "0" to the binary string.
			case uint8(1): // If the hash bit is 1.
				s += "1" // Append "1" to the binary string.
			default:
				panic("Hash value is not 0 or 1") // Panic if a hash bit is neither 0 nor 1 (should not happen).
			}
		}
		v, err := strconv.ParseUint(s, 2, 64) // Parse the binary string 's' into a uint64.
		if err != nil {
			panic(err) // Panic if there's an error during parsing (should not happen for valid binary strings).
		}
		basicKeys[i] = v // Assign the parsed uint64 value to the slice of basic keys.
	}
	return basicKeys // Return the slice of uint64 keys.
}
