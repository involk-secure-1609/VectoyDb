package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"vectorDb/client"
	"vectorDb/db"
	"vectorDb/store"

	"github.com/spf13/cobra"
)

var vectorDb *db.Db
var vectorClient *client.GeminiClient
var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveMode()
	},
}

// Available commands map
var commands = map[string]func([]string){
	"help": func(args []string) {
		fmt.Println("Available commands:")
		fmt.Println(" help    - Show this help")
		fmt.Println(" use storeName indexType -creates a database with either lsh or hnsw as the underlying data structure")
		fmt.Println(" search storeName key -searches for the key in the database")
		fmt.Println(" save storename -saves the database to disk")
		fmt.Println("  exit    - Exit the application")

		fmt.Println("  version - Show version information")
		fmt.Println("  exit    - Exit the application")
	},
	"version": func(args []string) {
		fmt.Println("VectoyDb version 1.0")
	},
	"exit": func(args []string) {
		fmt.Println("Goodbye!")
		os.Exit(0)
	},
	"use": func(args []string) {
		switch strings.ToLower(args[1]) {
		case "lsh":
			lshStore, err := store.NewLshStore()
			lshStore.Load(strings.ToLower(args[0]))
			if err != nil {
				log.Println(err)
				return
			}
			vectorDb.Store = lshStore
		case "hnsw":
			hnswStore, err := store.NewHnswStore()
			hnswStore.Load(strings.ToLower(args[0]))
			if err != nil {
				log.Println(err)
				return
			}
			vectorDb.Store = hnswStore
		default:
			log.Printf("store of type %s not availible", args[0])
		}
	},
	"save": func(args []string) {
		err := vectorDb.Store.Save(strings.ToLower(args[0]))
		if err != nil {
			log.Println(err)
			return
		}
	},
	"insert": func(args []string) {
		storeName := strings.ToLower(args[0])
		for _, key := range args[1:] {
			err := vectorDb.Insert(storeName, key)
			if err != nil {
				log.Printf("could not insert key :%s", key)
				continue
			}
		}
	},
	"delete": func(args []string) {
		storeName := strings.ToLower(args[0])
		for _, key := range args[1:] {
			_, err := vectorDb.Delete(storeName, key)
			if err != nil {
				log.Printf("could not insert key :%s", key)
				continue
			}
		}
	},
	"search": func(args []string) {
		storeName := strings.ToLower(args[0])

		limit, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("could not convert key to int: %v", err)
		}

		query := strings.Join(args[2:], " ")

		results, err := vectorDb.Search(storeName, query,limit)
		if err!=nil{
			log.Println(err)
			return
		}
		for result:=range(results){
			log.Println(result)
		}
		
	},
}

func runInteractiveMode() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to VectoyDb")
	fmt.Println("Type 'help' for available commands or 'exit' to quit.")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		// Trim whitespace and convert to lowercase
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Split the input into command and arguments
		parts := strings.Fields(input)
		command := strings.ToLower(parts[0])
		args := parts[1:]

		// Execute the command if it exists
		if cmd, exists := commands[command]; exists {
			cmd(args)
		} else {
			fmt.Printf("Unknown command: %s\nType 'help' for available commands.\n", command)
		}
	}
}

func Execute() {
	vectorClient = client.NewGeminiClient()
	vectorDb = db.NewVectorDb(vectorClient)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
