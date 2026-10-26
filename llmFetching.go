package main

import (
	"log"
	"os"
)

func getLLMres(prompt string) {
	apiKey := os.Getenv("GEMINI_API_KEY")

	//rollback ketika env tidak ada
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
		return
	}

	return
}
