package tests

import (
	// "os"
	"os"
	"testing"
	"vectorDb/store"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a test store with a named index
func TestHnswStoreSetup(t *testing.T)  {
	_, err := store.NewHnswStore()
	assert.NoError(t, err)
}


// Tests the Insert and Lookup functionality
func TestHnswStoreInsertLookup(t *testing.T) {
	storeName := "test_store"
	hnswStore, err := store.NewHnswStore()
	assert.NoError(t, err)

	// Insert embeddings for each letter
	for i := 'a'; i <= 'z'; i++ {
		embedding := generateRandomFloat32Array(8)
		err := hnswStore.Insert(storeName, embedding, string(i))
		assert.NoError(t, err)
	}


	embedding:=make([]float32,0)
	// Lookup each embedding
	for i := 'a'; i <= 'z'; i++ {
		embedding, err := hnswStore.Lookup(storeName,embedding ,string(i))
		assert.NoError(t, err)
		assert.Equal(t, 8, len(embedding))
	}
	
	// Test lookup for non-existent key
	_, err = hnswStore.Lookup(storeName, embedding,"non_existent_key")
	assert.Error(t, err)
	
}

// Tests the Insert and Delete functionality
func TestHnswStoreInsertDelete(t *testing.T) {
	storeName := "test_store"
	hnswStore, err := store.NewHnswStore()
	assert.NoError(t, err)
	
	// Insert a single embedding
	key := "a"
	embedding := generateRandomFloat32Array(8)
	err = hnswStore.Insert(storeName, embedding, key)
	assert.NoError(t, err)
	
	// Verify it can be looked up
	_, err = hnswStore.Lookup(storeName, embedding,key)
	assert.NoError(t, err)
	
	// Delete the embedding
	deleted, err := hnswStore.Delete(storeName, embedding, key)
	assert.NoError(t, err)
	assert.True(t, deleted)
	
	// Verify it's been deleted
	_, err = hnswStore.Lookup(storeName, embedding,key)
	assert.Error(t, err)
	
	// Try deleting again
	deleted, err = hnswStore.Delete(storeName, embedding, key)
	assert.NoError(t, err)
	assert.False(t, deleted)
}

// Tests the Search functionality
func TestHnswStoreSearch(t *testing.T) {
	storeName := "test_store"
	hnswStore, err := store.NewHnswStore()
	assert.NoError(t, err)

	// Insert multiple embeddings
	for i := 'a'; i <= 'z'; i++ {
		embedding := generateRandomFloat32Array(8)
		err := hnswStore.Insert(storeName, embedding, string(i))
		assert.NoError(t, err)
	}
	
	// Search for nearest neighbors
	query := generateRandomFloat32Array(8)
	results, err := hnswStore.Search(storeName, query, 5)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(results), 5)
	
	// With fewer items than limit
	smallStoreName := "small_store"
	
	err = hnswStore.Insert(smallStoreName, generateRandomFloat32Array(8), "a")
	assert.NoError(t, err)
	
	results, err = hnswStore.Search(smallStoreName, query, 5)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
}

// Tests the Save and Load functionality
func TestHnswStoreSaveLoad(t *testing.T) {
	storeName := "test_store"
	hnswStore, err := store.NewHnswStore()
	assert.NoError(t, err)
	defer os.Remove(storeName + "_hnsw" + ".store")
	
	// Insert embeddings
	for i := 'a'; i <= 'e'; i++ {
		embedding := generateRandomFloat32Array(8)
		err := hnswStore.Insert(storeName, embedding, string(i))
		assert.NoError(t, err)
	}
	
	// Save the store
	err = hnswStore.Save(storeName)
	assert.NoError(t, err)
	embedding := make([]float32,0)
	// Clear the store by deleting all entries
	for i := 'a'; i <= 'e'; i++ {
		embedding, err := hnswStore.Lookup(storeName,embedding, string(i))
		assert.NoError(t, err)
		deleted, err := hnswStore.Delete(storeName, embedding, string(i))
		assert.NoError(t, err)
		assert.True(t, deleted)
	}
	// Verify entries are gone
	_, err = hnswStore.Lookup(storeName, embedding,"a")
	assert.Error(t, err)
	
	// Load the store
	err = hnswStore.Load(storeName)
	assert.NoError(t, err)
	
	// Verify entries are restored
	for i := 'a'; i <= 'e'; i++ {
		embedding, err := hnswStore.Lookup(storeName,embedding, string(i))
		assert.NoError(t, err)
		assert.Equal(t, 8, len(embedding))
	}
}

// Tests multiple named stores within the same HnswStore
func TestHnswStoreMultipleStores(t *testing.T) {
	store1 := "store1"
	store2 := "store2"
	hnswStore, err := store.NewHnswStore()
	assert.NoError(t, err)

	// Insert different embeddings in different stores
	for i := 'a'; i <= 'e'; i++ {
		embedding := generateRandomFloat32Array(8)
		err := hnswStore.Insert(store1, embedding, string(i))
		assert.NoError(t, err)
	}
	
	for i := 'v'; i <= 'z'; i++ {
		embedding := generateRandomFloat32Array(8)
		err := hnswStore.Insert(store2, embedding, string(i))
		assert.NoError(t, err)
	}
	
	embedding := make([]float32,0)
	// Verify lookups work correctly
	_, err = hnswStore.Lookup(store1,embedding, "a")
	assert.NoError(t, err)
	
	_, err = hnswStore.Lookup(store2,embedding, "z")
	assert.NoError(t, err)
	
	// Cross-store lookups should fail
	_, err = hnswStore.Lookup(store1,embedding, "z")
	assert.Error(t, err)
	
	_, err = hnswStore.Lookup(store2,embedding, "a")
	assert.Error(t, err)
}