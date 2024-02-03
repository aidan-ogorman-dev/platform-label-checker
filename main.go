package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello from your go program")
	ghRefName := os.Getenv("GITHUB_REF_NAME")
	ghRefId := os.Getenv("GITHUB_REF")
	fmt.Printf("Running against ref: %s with ref ID: %s", ghRefName, ghRefId)
}