package store

type Store interface{
	Insert(storeName string,embedding []float32,key string) (error)
	Search(storeName string,query []float32, limit int) ([]string,error)
	Delete(storeName string,embedding []float32,key string) (bool,error)
	Load(storeName string) (error)
	Save(storeName string) (error)
	Lookup(storeName string,embedding []float32,key string) ([]float32,error)
}