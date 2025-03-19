package tests

import (
	"os"
	"testing"
	"vectorDb/hnsw"

	"github.com/stretchr/testify/assert"
)


// Tests the Dims and Len function
func TestHNSWInsertDimLen(t *testing.T) {
	hnswGraph := hnsw.NewHNSWGraph[string]("")
	embedding := generateRandomFloat64Array(8)
	node := hnsw.MakeNode[string]("a", embedding)

	hnswGraph.Insert(node)

	size := hnswGraph.Len()
	dim := hnswGraph.Dims()
	assert.Equal(t, 1, size)
	assert.Equal(t, 8, dim)
	assert.NotEqual(t, 7, dim)

}

// Tests the Insert and Lookup functionality
func TestHNSWInsertLookup(t *testing.T) {
	hnswGraph := hnsw.NewHNSWGraph[string]("")
	for i := 'a'; i <= 'z'; i++ {
		embedding := generateRandomFloat64Array(8)
		node := hnsw.MakeNode[string](string(i), embedding)
		hnswGraph.Insert(node)
	}

	for i := 'a'; i <= 'z'; i++ {
		embedding, present := hnswGraph.Lookup(string(i))
		assert.Equal(t, 8, len(embedding))
		assert.NotEqual(t, 9, len(embedding))
		assert.Equal(t, true, present)
		assert.NotEqual(t, false, present)
	}
}

// Tests the Insert and Delete functionality
func TestHNSWInsertDelete(t *testing.T) {
	hnswGraph := hnsw.NewHNSWGraph[string]("")
	embedding := generateRandomFloat64Array(8)
	node := hnsw.MakeNode[string]("a", embedding)

	deleted := hnswGraph.Delete(node.Key)
	assert.Equal(t, false, deleted)

	hnswGraph.Insert(node)
	length := hnswGraph.Len()
	assert.Equal(t, 1, length)

	deleted = hnswGraph.Delete(node.Key)
	assert.Equal(t, true, deleted)
	length = hnswGraph.Len()
	assert.Equal(t, 0, length)

	deleted = hnswGraph.Delete(node.Key)
	assert.Equal(t, false, deleted)
	length = hnswGraph.Len()
	assert.Equal(t, 0, length)
}

// Tests the Load Functionality
func TestHNSWLoad(t *testing.T) {
	hnswGraph := hnsw.NewHNSWGraph[string]("")
	testFile := "test"
	defer os.Remove(testFile + "_hnsw" + ".store")
	err := hnswGraph.Load(testFile)
	assert.Equal(t, nil, err)
}

// Tests the Insert and Save Functionality
func TestHNSWInsertSave(t *testing.T) {
	hnswGraph := hnsw.NewHNSWGraph[string]("")
	testFile := "test"
	defer os.Remove(testFile + "_hnsw" + ".store")

	embedding := generateRandomFloat64Array(8)
	node := hnsw.MakeNode[string]("a", embedding)
	hnswGraph.Insert(node)
	err := hnswGraph.Save(testFile)
	assert.Equal(t, nil, err)
}

// Tests the Save and Load Functionality
func TestHNSWSaveLoad_1(t *testing.T) {
	hnswGraph := hnsw.NewHNSWGraph[string]("")
	testFile := "test"
	defer os.Remove(testFile + "_hnsw" + ".store")

	embedding := generateRandomFloat64Array(8)
	// log.Println(embedding)
	node := hnsw.MakeNode[string]("a", embedding)
	hnswGraph.Insert(node)
	err := hnswGraph.Save(testFile)
	assert.Equal(t, nil, err)
	deleted := hnswGraph.Delete(node.Key)
	assert.Equal(t, true, deleted)
	embedding, present := hnswGraph.Lookup(node.Key)
	assert.Equal(t, false, present)
	err = hnswGraph.Load(testFile)
	assert.Equal(t, nil, err)

	embedding, present = hnswGraph.Lookup(node.Key)
	assert.Equal(t, 8, len(embedding))
	assert.Equal(t, true, present)
}

// Tests the Save and Load Functionality
func TestHNSWSaveLoad_2(t *testing.T) {
	hnswGraph := hnsw.NewHNSWGraph[string]("")
	testFile := "test"
	defer os.Remove(testFile + "_hnsw" + ".store")

	for i := 'a'; i <= 'z'; i++ {
		embedding := generateRandomFloat64Array(8)
		node := hnsw.MakeNode[string](string(i), embedding)
		hnswGraph.Insert(node)
	}
	lenBeforeSaving:=hnswGraph.Len()
	err:=hnswGraph.Save(testFile)
	assert.Equal(t,nil,err)

	for i := 'a'; i <= 'z'; i++ {
		hnswGraph.Delete(string(i))
	}

	lenAfterDeleting:=hnswGraph.Len()
	assert.Equal(t,0,lenAfterDeleting)

	err=hnswGraph.Load(testFile)
	assert.Equal(t,nil,err)

	lenAfterLoading:=hnswGraph.Len()
	assert.Equal(t,lenBeforeSaving,lenAfterLoading)
}
