package store

import (
	"fmt"
	"vectorDb/hnsw"
)

// Contains all the store functionality of the VectorDb

// type Store interface{
// 	Insert()
// 	Delete()
// 	Search()
// 	Save()
// 	Load()
// }

type HnswStore struct {
	store map[string]*hnsw.HNSWGraph[string]
}

func NewHnswStore() (Store, error) {
	hnswStore := &HnswStore{
		store: map[string]*hnsw.HNSWGraph[string]{},
	}

	return hnswStore, nil
}


func (hnswStore *HnswStore) Search(storeName string,query []float64, limit int) ([]string,error) {
	neighborNodes:=hnswStore.store[storeName].Search(query,limit)
	neighbors:=make([]string,0)
	for _,hnswNode:=range(neighborNodes){
		neighbors=append(neighbors, hnswNode.Key)
	}

	return neighbors,nil
}

func (hnswStore *HnswStore) Insert(storeName string,embedding []float64,key string) (error){
	hnswStore.store[storeName].Insert(hnsw.Node[string]{Key: key,Embed: embedding})
	return nil
}


func (hnswStore *HnswStore) Lookup(storeName string,key string) ([]float64,error) {
	embedding,present:=hnswStore.store[storeName].Lookup(key)
	if !present{
		return nil,fmt.Errorf("key not present in the database")
	}
	return embedding,nil
}

func (hnswStore *HnswStore) Delete(storeName string,embdedding []float64,key string) (bool,error) {
	deleted:=hnswStore.store[storeName].Delete(key);
	return deleted,nil
}

func (hnswStore *HnswStore) Load(storeName string) (error) {
	err:=hnswStore.store[storeName].Load(storeName);
	return err
}

func (hnswStore *HnswStore) Save(storeName string) (error) {
	err:=hnswStore.store[storeName].Save(storeName);
	return err
}


// func (hnswStore *VectorStore) save(directoryName string) (error) {
// 	err:=hnswStore.directoryStore[directoryName].save()
// 	return err
// }

