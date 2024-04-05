package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-be-ai/server"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}
	port := os.Getenv("PORT")
	fmt.Printf("Server is listening on port %s...\n", port)
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/play/", server.PlayGameHandler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error:", err)
	}
}
