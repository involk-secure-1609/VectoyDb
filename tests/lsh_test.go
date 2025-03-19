package tests

import (
	"os"
	"testing"
	"vectorDb/lsh"

	"github.com/stretchr/testify/assert"
)

func TestLSHInsert(t *testing.T) {
	lshIndex:=lsh.NewCosineLsh(20,15,15,"euclidean")
	embedding := generateRandomFloat32Array(20)
	lshIndex.Insert(embedding,"a")
	present:=lshIndex.Lookup(embedding,"a")
	assert.Equal(t,true,present)
}

func TestLSHDelete(t *testing.T) {
	lshIndex:=lsh.NewCosineLsh(20,15,15,"euclidean")
	embedding := generateRandomFloat32Array(20)
	lshIndex.Insert(embedding,"a")
	present:=lshIndex.Lookup(embedding,"a")
	assert.Equal(t,true,present)

	lshIndex.Delete(embedding,"a")
	present=lshIndex.Lookup(embedding,"a")
	assert.Equal(t,false,present)
}

func TestLSHLoad(t *testing.T) {
	lshIndex:=lsh.NewCosineLsh(20,15,15,"euclidean")
	testFile := "test"
	defer os.Remove(testFile + "_lsh" + ".store")
	err:=lshIndex.Load(testFile)
	assert.Equal(t,nil,err)
}

func TestLSHSave(t *testing.T) {
	lshIndex:=lsh.NewCosineLsh(20,15,15,"euclidean")
	testFile := "test"
	defer os.Remove(testFile + "_lsh" + ".store")
	err:=lshIndex.Load(testFile)
	assert.Equal(t,nil,err)
	embedding := generateRandomFloat32Array(20)
	lshIndex.Insert(embedding,"a")
	err=lshIndex.Save(testFile)
	assert.Equal(t,nil,err)
}

func TestLSHSaveAndThenLoad1 (t *testing.T) {
	lshIndex:=lsh.NewCosineLsh(20,15,15,"euclidean")
	testFile := "test"
	defer os.Remove(testFile + "_lsh" + ".store")
	err:=lshIndex.Load(testFile)
	assert.Equal(t,nil,err)
	embeddings:=make([][]float32,0)
	for i := 'a'; i <= 'z'; i++ {
		embedding := generateRandomFloat32Array(8)
		embeddings = append(embeddings,embedding )
		lshIndex.Insert(embedding,string(i))
	}
	err=lshIndex.Save(testFile)
	assert.Equal(t,nil,err)

	for i := 'a'; i <= 'z'; i++ {
		index := int(i - 'a')
		lshIndex.Delete(embeddings[int(index)],string(i))
	}
	var present bool
	for i := 'a'; i <= 'z'; i++ {
		index := int(i - 'a')
		present=lshIndex.Lookup(embeddings[int(index)],string(i))
		assert.Equal(t,false,present)
	}

	err=lshIndex.Load(testFile)
	assert.Equal(t,nil,err)

	for i := 'a'; i <= 'z'; i++ {
		index := int(i - 'a')
		present=lshIndex.Lookup(embeddings[int(index)],string(i))
		assert.Equal(t,true,present)
	}
}



func TestLSHSaveAndThenLoad2 (t *testing.T) {
	lshIndex:=lsh.NewCosineLsh(20,15,15,"euclidean")
	testFile := "test"
	defer os.Remove(testFile + "_lsh" + ".store")
	err:=lshIndex.Load(testFile)
	assert.Equal(t,nil,err)
	embedding := generateRandomFloat32Array(20)
	lshIndex.Insert(embedding,"a")
	err=lshIndex.Save(testFile)
	assert.Equal(t,nil,err)

	lshIndex.Delete(embedding,"a")
	present:=lshIndex.Lookup(embedding,"a")
	assert.Equal(t,false,present)

	err=lshIndex.Load(testFile)
	assert.Equal(t,nil,err)

	present=lshIndex.Lookup(embedding,"a")
	assert.Equal(t,true,present)
}
