package client

type Client interface{
	Embed(key string) ([]float32, error)
}
