package store

type Store interface{
	Insert(storeName string,embedding []float64,key string) (error)
	Search(storeName string,query []float64, limit int) ([]string,error)
	Delete(storeName string,embedding []float64,key string) (bool,error)
	Load(storeName string) (error)
	Save(storeName string) (error)
}