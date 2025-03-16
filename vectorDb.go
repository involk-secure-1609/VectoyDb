package main

import "vectorDb/store"

type Db struct {
	Client *GeminiClient
	Store  store.Store
}

func NewVectorDb(client *GeminiClient, store store.Store) *Db {
	return &Db{Client: client, Store: store}
}


func (db *Db) Search(directory string, query string, limit int) ([]string, error) {
	embedding, err := db.Client.embed(query)
	if err != nil {
		return []string{}, err
	}
	if limit == -1 {
		limit = 3
	}
	queryResult,err := db.Store.Search(directory,Float32ToFloat64(embedding),limit)
	if err!=nil{
		return nil,err
	}

	return queryResult, nil

}

func (db *Db) Insert(directory string, key string) error {
	embedding, err := db.Client.embed(key)
	if err != nil {
		return err
	}
	err = db.Store.Insert(directory,Float32ToFloat64(embedding),key)
	return err
}

// func (db *Db) Loo(directory string, key string) ([]float32, error) {
// 	embedding, err := db.vectorStore.lookup(directory, key)
// 	if err != nil {
// 		return []float32{}, err
// 	}
// 	return embedding, nil
// }

func (db *Db) Delete(directory string, key string) (bool, error) {
	embedding, err := db.Client.embed(key)
	if err != nil {
		return false,err
	}
	deleted,err := db.Store.Delete(directory,Float32ToFloat64(embedding),key)
	return deleted,err
}

