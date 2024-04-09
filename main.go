package main

import (
	"fmt"
	"go-be-ai/server"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}
}
func main() {
	port := os.Getenv("PORT")
	fmt.Printf("Server is listening on port %s...\n", port)
	http.HandleFunc("/ping", server.PingHandler)
	http.HandleFunc("/play/", server.PlayGameHandler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error:", err)
	}
}
