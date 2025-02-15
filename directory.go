package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coder/hnsw"
	"github.com/google/renameio"
)


type Directory struct{
	directoryFile *os.File
	directoryName string
	graph *hnsw.Graph[string]
}


func NewDirectory(directoryName string) *Directory{
	return &Directory{
		directoryName: directoryName+".graph",
	}
}
func (d *Directory) init() error{
	directoryFile, err := os.OpenFile(d.directoryName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil{
		return err
	}
	defer directoryFile.Close()
	info,err:=directoryFile.Stat()
	if err!=nil{
		return err
	}
	if (info.Size()>0){
		err = d.graph.Import(bufio.NewReader(directoryFile))
		if err != nil {
			return fmt.Errorf("import: %w", err)
		}
	}

	go func(){
		time.Sleep(1*time.Minute)
		d.save()
	}()
	d.directoryFile = directoryFile
	return nil
}

func (d *Directory) save() error{
	tmp, err := renameio.TempFile("", d.directoryFile.Name())
	if err != nil {
		return err
	}
	defer tmp.Cleanup()

	wr := bufio.NewWriter(tmp)
	err = d.graph.Export(wr)
	if err != nil {
		return fmt.Errorf("exporting: %w", err)
	}

	err = wr.Flush()
	if err != nil {
		return fmt.Errorf("flushing: %w", err)
	}

	err = tmp.CloseAtomicallyReplace()
	if err != nil {
		return fmt.Errorf("closing atomically: %w", err)
	}

	return nil
}
func (d *Directory) insert(key string,embedding []float32) {
	d.graph.Add(hnsw.MakeNode(key,embedding))
}

func (d *Directory) query(query []float32,limit int) []hnsw.Node[string]{
	neighbours:=d.graph.Search(query, limit)
	for _,neighbor:=range(neighbours){
		log.Print(neighbor.Key)
		log.Print(":")
		log.Println(neighbor.Value)
	}

	return neighbours
}


func (d *Directory) lookup(key string) (hnsw.Vector){
	vector,ok:=d.graph.Lookup(key)
	if !ok{
		return nil
	}
	return vector
}


func (d *Directory) delete(key string) (bool){
	ok:=d.graph.Delete(key)
	return ok
}
