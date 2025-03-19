package db

import (
	"vectorDb/client"
	"vectorDb/store"
)

type Db struct {
	Client client.Client
	Store  store.Store
}

func NewVectorDbWithClient(client client.Client) *Db {
	return &Db{Client: client}
}

func NewVectorDbWithStore(store store.Store) *Db {
	return &Db{Store: store}
}

func NewVectorDbWithClientAndStore(client client.Client,store store.Store) *Db {
	return &Db{Client:client,Store: store}
}


func (db *Db) Search(storeName string, query string, limit int) ([]string, error) {
	embedding, err := db.Client.Embed(query)
	if err != nil {
		return nil, err
	}
	if limit == -1 {
		limit = 3
	}
	queryResult,err := db.Store.Search(storeName,(embedding),limit)
	if err!=nil{
		return nil,err
	}

	return queryResult, nil

}

func (db *Db) Insert(storeName string, key string) error {
	embedding, err := db.Client.Embed(key)
	if err != nil {
		return err
	}
	err = db.Store.Insert(storeName,embedding,key)
	return err
}

func (db *Db) Lookup(storeName string, key string) ([]float32, error) {
	embedding, err := db.Client.Embed(key)
	if err != nil {
		return nil, err
	}
	searchResult, err := db.Store.Lookup(storeName,embedding,key)
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (db *Db) Delete(storeName string, key string) (bool, error) {
	embedding, err := db.Client.Embed(key)
	if err != nil {
		return false,err
	}
	deleted,err := db.Store.Delete(storeName,(embedding),key)
	return deleted,err
}

func (db *Db) Save(storeName string) (error) {
	err:=db.Store.Save(storeName)
	return err
}

func (db *Db) Load(storeName string) (error) {
	err := db.Store.Load(storeName)
	return err
}

