package lsh

// import "math/rand"

// type Point struct {
// 	key       string
// 	embedding []float32
// }
// type HashKey string
// type HashBucket []Point
// type HashTable []map[uint64]HashBucket
// type DistanceFunc func(Point, Point) float32

// func NewHyperplanes(h,d int) [][]float32{
// 	hyperplanes := make([][]float32,h)
// 	for i :=range(h){
// 		hyperplane:=[]float32{}
// 		for j :=range(d){
// 			point:=rand.NormFloat64()
// 		}
// 	}
// }
// type Lsh struct {
// 	d            int          //dimension of the embeddings and thereby the hyperplanes
// 	hyperplanes  [][]float32  // hyperplanes being used
// 	distanceFunc DistanceFunc // Function to calculate the distance between vectors.
// 	l            int          // Number of hash tables to use.
// 	m            int          // Number of hyperplanes used in each hash table to create a hash key.
// 	h            int          // Total number of hyperplanes (l * m).
// 	hashTable *HashTable      // Table used to store all the hashes
// }

// func NewLsh(d int, l int,m int){

// 	lsh := &Lsh{
// 		d:d,
// 		l:l,
// 		m:m,
// 	}

// 	lsh.hyperplanes=NewHyperplanes()
// }
