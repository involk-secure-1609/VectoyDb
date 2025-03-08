package main

import (
	"log"
	"net"
)

func main() {
	vectorClient := NewVectorClient("localhost:11434", "mistral")
	err:=vectorClient.start()
	if err!=nil{
		log.Println(err)
		log.Println("To start the server, you can run:")
		log.Printf("ollama serve")	
		return
	}

	vectorStore,err:=NewVectorStore()
	if err!=nil{
		log.Println(err)
		return 
	}
	vectorDb:=NewVectorDb(vectorClient,vectorStore);
	log.Println(vectorDb.vectorClient.model)
	l,err:=net.Listen("tcp","localhost:8050")
	if err!=nil{
		log.Println(err)
		return  
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {
		msg:=make([]byte,1024)
		n,err:=conn.Read(msg)
		if err!=nil{
			log.Println(err)
			break
		}
		msg=msg[:n]
		log.Println(string(msg))
	}
	
}
