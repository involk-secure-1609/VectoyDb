package store

import (
	"errors"
	"vectorDb/lsh"
)

type LshStore struct {
	store map[string]*lsh.CosineLsh
}

func NewLshStore() (Store, error) {
	lshStore := &LshStore{
		store:map[string]*lsh.CosineLsh{},
	}
	return lshStore, nil
}

func (lshStore *LshStore) initialize(storeName string){
	_,present:=lshStore.store[storeName]
	if !present{
		lshStore.store[storeName]=lsh.NewCosineLsh(20,15,15,"euclidean")
	}
}

func (lshStore *LshStore) Search(storeName string,query []float64, limit int) ([]string,error) {
	lshStore.initialize(storeName)
	searchResults:=lshStore.store[storeName].Search(query,limit)
	results:=make([]string,0)
	for _,result:=range(searchResults){
		results=append(results, result.ExtraData)
	}

	return results,nil
}

func (lshStore *LshStore) Insert(storeName string,embedding []float64,key string) (error){
	lshStore.initialize(storeName)
	lshStore.store[storeName].Insert(embedding,key)
	return nil
}


func (lshStore *LshStore) Lookup(storeName string,embedding []float64,key string) ([]float64,error) {
	lshStore.initialize(storeName)
	present:=lshStore.store[storeName].Lookup(embedding,key)
	if !present{
		return nil,errors.New("key not present in the database")
	}
	return embedding,nil
}

func (lshStore *LshStore) Delete(storeName string,embdedding []float64,key string) (bool,error) {
	lshStore.initialize(storeName)
	lshStore.store[storeName].Delete(embdedding,key);
	return true,nil
}

func (lshStore *LshStore) Load(storeName string) (error) {
	lshStore.initialize(storeName)
	err:=lshStore.store[storeName].Load(storeName);
	return err
}

func (lshStore *LshStore) Save(storeName string) (error) {
	lshStore.initialize(storeName)
	err:=lshStore.store[storeName].Save(storeName);
	return err
}


