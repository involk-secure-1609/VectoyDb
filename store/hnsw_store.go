package store

import (
	"fmt"
	"vectorDb/hnsw"
)

type HnswStore struct {
	store map[string]*hnsw.HNSWGraph[string]
}

func NewHnswStore() (Store, error) {
	hnswStore := &HnswStore{
		store: map[string]*hnsw.HNSWGraph[string]{},
	}

	return hnswStore, nil
}

func (hnswStore *HnswStore) initialize(storeName string){
	_,present:=hnswStore.store[storeName]
	if !present{
		hnswStore.store[storeName]=hnsw.NewHNSWGraph[string]("")
	}
}
func (hnswStore *HnswStore) Search(storeName string,query []float64, limit int) ([]string,error) {
	hnswStore.initialize(storeName)
	neighborNodes:=hnswStore.store[storeName].Search(query,limit)
	neighbors:=make([]string,0)
	for _,hnswNode:=range(neighborNodes){
		neighbors=append(neighbors, hnswNode.Key)
	}

	return neighbors,nil
}

func (hnswStore *HnswStore) Insert(storeName string,embedding []float64,key string) (error){
	hnswStore.initialize(storeName)
	hnswStore.store[storeName].Insert(hnsw.Node[string]{Key: key,Embed: embedding})
	return nil
}


func (hnswStore *HnswStore) Lookup(storeName string,embedding []float64,key string) ([]float64,error) {
	hnswStore.initialize(storeName)
	embeddingFound,present:=hnswStore.store[storeName].Lookup(key)
	if !present{
		return nil,fmt.Errorf("key not present in the database")
	}
	return embeddingFound,nil
}

func (hnswStore *HnswStore) Delete(storeName string,embdedding []float64,key string) (bool,error) {
	hnswStore.initialize(storeName)
	deleted:=hnswStore.store[storeName].Delete(key);
	return deleted,nil
}

func (hnswStore *HnswStore) Load(storeName string) (error) {
	hnswStore.initialize(storeName)
	err:=hnswStore.store[storeName].Load(storeName);
	return err
}

func (hnswStore *HnswStore) Save(storeName string) (error) {
	hnswStore.initialize(storeName)
	err:=hnswStore.store[storeName].Save(storeName);
	return err
}


