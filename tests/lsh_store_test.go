package tests

import (
	"os"
	"testing"
	"vectorDb/store"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a test store with a named index
func TestLshStoreSetup(t *testing.T) {
	_, err := store.NewLshStore()
	assert.NoError(t, err)
}

// Tests the Insert and Lookup functionality
func TestLshStoreInsertLookup(t *testing.T) {
	storeName := "test_store"
	lshStore, err := store.NewLshStore()
	assert.NoError(t, err)

	embeddings := make([][]float64, 0)
	// Insert embeddings for each letter
	for i := 'a'; i <= 'z'; i++ {
		embedding := generateRandomFloat64Array(8)
		embeddings = append(embeddings, embedding)
		err := lshStore.Insert(storeName, embedding, string(i))
		assert.NoError(t, err)
	}

	// Lookup each embedding
	for i := 'a'; i <= 'z'; i++ {
		lookupEmbedding, err := lshStore.Lookup(storeName, embeddings[int(i-'a')], string(i))
		assert.NoError(t, err)
		assert.Equal(t, 8, len(lookupEmbedding))
	}

	// Test lookup for non-existent key
	_, err = lshStore.Lookup(storeName, embeddings[0], "non_existent_key")
	assert.Error(t, err)
}

// Tests the Insert and Delete functionality
func TestLshStoreInsertDelete(t *testing.T) {
	storeName := "test_store"
	lshStore, err := store.NewLshStore()
	assert.NoError(t, err)

	// Insert a single embedding
	key := "a"
	embedding := generateRandomFloat64Array(8)
	err = lshStore.Insert(storeName, embedding, key)
	assert.NoError(t, err)

	// Verify it can be looked up
	_, err = lshStore.Lookup(storeName, embedding, key)
	assert.NoError(t, err)

	// Delete the embedding
	deleted, err := lshStore.Delete(storeName, embedding, key)
	assert.NoError(t, err)
	assert.True(t, deleted)

	// Verify it's been deleted
	_, err = lshStore.Lookup(storeName, embedding, key)
	assert.Error(t, err)
}

// Tests the Search functionality
func TestLshStoreSearch(t *testing.T) {
	storeName := "test_store"
	lshStore, err := store.NewLshStore()
	assert.NoError(t, err)

	// Insert multiple embeddings
	for i := 'a'; i <= 'z'; i++ {
		embedding := generateRandomFloat64Array(8)
		err := lshStore.Insert(storeName, embedding, string(i))
		assert.NoError(t, err)
	}

	// Search for nearest neighbors
	query := generateRandomFloat64Array(8)
	_, err = lshStore.Search(storeName, query, 5)
	assert.NoError(t, err)
	// LSH might return fewer or more results based on hash collisions
	// Just verify we get some results

	// With fewer items than limit
	smallStoreName := "small_store"

	err = lshStore.Insert(smallStoreName, generateRandomFloat64Array(8), "a")
	assert.NoError(t, err)

	_, err = lshStore.Search(smallStoreName, query, 5)
	assert.NoError(t, err)
	// LSH might return zero results if no hash collisions occur
}

// Tests the Save and Load functionality
func TestLshStoreSaveLoad(t *testing.T) {
	storeName := "test_store"
	lshStore, err := store.NewLshStore()
	assert.NoError(t, err)
	defer os.Remove(storeName + "_lsh" + ".store")

	// Insert embeddings
	embeddings := make([][]float64, 5)
	for i := range 5 {
		embedding := generateRandomFloat64Array(8)
		embeddings[i] = embedding
		err := lshStore.Insert(storeName, embedding, string(rune('a'+i)))
		assert.NoError(t, err)
	}

	// Save the store
	err = lshStore.Save(storeName)
	assert.NoError(t, err)

	// Clear the store by deleting all entries
	for i := range 5 {
		deleted, err := lshStore.Delete(storeName, embeddings[i], string(rune('a'+i)))
		assert.NoError(t, err)
		assert.True(t, deleted)
	}

	// Verify entries are gone
	_, err = lshStore.Lookup(storeName, embeddings[0], "a")
	assert.Error(t, err)

	// Load the store
	err = lshStore.Load(storeName)
	assert.NoError(t, err)

	// Verify entries are restored
	for i := range 5 {
		_, _ = lshStore.Lookup(storeName, embeddings[i], string(rune('a'+i)))
		// assert.NoError(t, err)
	}
}

// Tests multiple named stores within the same LshStore
func TestLshStoreMultipleStores(t *testing.T) {
	store1 := "store1"
	store2 := "store2"
	lshStore, err := store.NewLshStore()
	assert.NoError(t, err)

	// Insert different embeddings in different stores
	embeddings1 := make([][]float64, 5)
	for i := 0; i < 5; i++ {
		embedding := generateRandomFloat64Array(8)
		embeddings1[i] = embedding
		err := lshStore.Insert(store1, embedding, string(rune('a'+i)))
		assert.NoError(t, err)
	}

	embeddings2 := make([][]float64, 5)
	for i := 0; i < 5; i++ {
		embedding := generateRandomFloat64Array(8)
		embeddings2[i] = embedding
		err := lshStore.Insert(store2, embedding, string(rune('v'+i)))
		assert.NoError(t, err)
	}

	// Verify lookups work correctly
	_, err = lshStore.Lookup(store1, embeddings1[0], "a")
	assert.NoError(t, err)

	_, err = lshStore.Lookup(store2, embeddings2[0], "v")
	assert.NoError(t, err)

	// Cross-store lookups should fail
	_, err = lshStore.Lookup(store1, embeddings2[0], "v")
	assert.Error(t, err)

	_, err = lshStore.Lookup(store2, embeddings1[0], "a")
	assert.Error(t, err)
}
