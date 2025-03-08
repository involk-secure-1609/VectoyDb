package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Contains all the store functionality of the VectorDb

type VectorStore struct {
	directoryStore map[string]*Directory
}

func NewVectorStore() (*VectorStore, error) {
	vectorStore := &VectorStore{
		directoryStore: map[string]*Directory{},
	}
	directoryStoreFile, err := os.OpenFile("directoryStore.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer directoryStoreFile.Close()
	info, err := directoryStoreFile.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() > 0 {
		offset := 0
		directoryStoreFile.Seek(0, io.SeekStart)
		numberOfDirectories := make([]byte, 4)
		n, err := directoryStoreFile.ReadAt(numberOfDirectories, int64(offset))
		if err != nil {
			return nil, err
		}
		numOfDirectories := binary.BigEndian.Uint32(numberOfDirectories)
		offset += n
		for range int(numOfDirectories) {
			directoryNameLength := make([]byte, 4)
			n, err := directoryStoreFile.ReadAt(numberOfDirectories, int64(offset))
			if err != nil {
				return nil, err
			}
			offset += n
			DirectoryNameLength := binary.BigEndian.Uint32(directoryNameLength)
			directoryName := make([]byte, DirectoryNameLength)
			n, err = directoryStoreFile.ReadAt(directoryName, int64(offset))
			if err != nil {
				return nil, err
			}
			offset += n
			vectorStore.directoryStore[string(directoryName)] = NewDirectory(string(directoryName))
		}

		for _, directory := range vectorStore.directoryStore {
			directory.init()
		}

	}

	return vectorStore, nil
}

func (vc *VectorStore) createDirectory(directoryName string) error {
	if vc.directoryStore[directoryName] != nil {
		return fmt.Errorf("directory %s already exists", directoryName)
	}
	vc.directoryStore[directoryName] = NewDirectory(directoryName)
	err := vc.directoryStore[directoryName].init()
	if err != nil {
		return err
	}
	return nil
}

func (vc *VectorStore) query(directory string,query []float32, limit int) ([]string) {
	neighborNodes:=vc.directoryStore[directory].query(query, limit)
	neighbors:=make([]string,0)
	for _,hnswNode:=range(neighborNodes){
		neighbors=append(neighbors, hnswNode.Key)
	}

	return neighbors
}

func (vc *VectorStore) insert(directory string,embedding []float32,key string) (error){
	vc.directoryStore[directory].insert(key,embedding)
	return nil
}


func (vc *VectorStore) lookup(directoryName string,key string) ([]float32,error) {
	vector:=vc.directoryStore[directoryName].lookup(key);
	if vector==nil{
		return nil,fmt.Errorf("key not present in the database")
	}
	return vector,nil
}

func (vc *VectorStore) delete(directoryName string,key string) (bool,error) {
	deleted:=vc.directoryStore[directoryName].delete(key);
	if deleted==false{
		return deleted,fmt.Errorf("key not present in the database")
	}
	return deleted,nil
}

func (vc *VectorStore) save(directoryName string) (error) {
	err:=vc.directoryStore[directoryName].save()
	return err
}

