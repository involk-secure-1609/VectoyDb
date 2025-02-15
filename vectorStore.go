package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Contains all the store functionality of the VectorDb

type VectorStore struct {
	directoryStore map[string]*Directory
	ollamaClient   *OllamaClient
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

func (vc *VectorStore) query(directory string, query string, limit int) error {
	inputs := make([]string, 0)
	inputs = append(inputs, query)
	response, err := vc.ollamaClient.embed(inputs)
	if err != nil {
		return err
	}

	responseDecoder := json.NewDecoder(response.Body)
	var r EmbedResponsePayload
	err = responseDecoder.Decode(&r)
	if err != nil {
		log.Printf("JSON decoding error: %v", err)
		return fmt.Errorf("failed to decode JSON response: %v", err)
	}

	log.Println(r.Embeddings)

	vc.directoryStore[directory].query(r.Embeddings[0], limit)
	return nil
}

func (vc *VectorStore) embed(input string) ([][]float32, error) {
	inputs := make([]string, 0)
	inputs = append(inputs, input)

	response, err := vc.ollamaClient.embed(inputs)
	if err != nil {
		return nil, err
	}

	responseDecoder := json.NewDecoder(response.Body)
	var r EmbedResponsePayload
	err = responseDecoder.Decode(&r)
	if err != nil {
		log.Printf("JSON decoding error: %v", err)
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	log.Println(r.Embeddings)
	return r.Embeddings, err
}

func (vc *VectorStore) lookup(directoryName string,key string) ([]float32,error) {
	vector:=vc.directoryStore[directoryName].lookup(key);
	if vector==nil{
		return nil,fmt.Errorf("key not present in the database")
	}
	return vector,nil
}
