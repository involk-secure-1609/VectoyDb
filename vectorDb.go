package main

type VectorDb struct {
	vectorClient *VectorClient
	vectorStore  *VectorStore
}

func NewVectorDb(vectorClient *VectorClient, vectorStore *VectorStore) *VectorDb {
	return &VectorDb{vectorClient: vectorClient, vectorStore: vectorStore}
}

func (vectorDb *VectorDb) query(directory string, query string, limit int) ([]string, error) {
	embedding, err := embeddingResponseProcessor(vectorDb.vectorClient.embed([]string{query}))
	if err != nil {
		return []string{}, err
	}
	if limit == -1 {
		limit = 3
	}
	queryResult := vectorDb.vectorStore.query(directory, embedding, limit)

	return queryResult, nil

}

func (vectorDb *VectorDb) insert(directory string, key string) error {
	embedding, err := embeddingResponseProcessor(vectorDb.vectorClient.embed([]string{key}))
	if err != nil {
		return err
	}
	err = vectorDb.vectorStore.insert(directory, embedding, key)
	if err != nil {
		return err
	}
	return nil
}

func (vectorDb *VectorDb) lookup(directory string, key string) ([]float32, error) {
	embedding, err := vectorDb.vectorStore.lookup(directory, key)
	if err != nil {
		return []float32{}, err
	}
	return embedding, nil
}

func (vectorDb *VectorDb) delete(directory string, key string) (bool, error) {
	deleted, err := vectorDb.vectorStore.delete(directory, key)
	if err != nil {
		return deleted, err
	}
	return deleted, nil
}

// func (vectorDb *VectorDb) save(directory string) error {
// 	return vectorDb.vectorStore.save(directory)
// }
