package main

import (
	"fmt"
	"log"

	"github.com/lattots/embego/pkg/embed"
	"github.com/lattots/embego/pkg/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load Word2Vec embedding model
	w2v, err := util.LoadEmbeddingModel("../data/w2v.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Load SQLite database
	db, err := gorm.Open(sqlite.Open("../data/token_frequency.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Create a DefaultSentenceTokenizer
	sTokenizer, err := util.LoadSentenceTokenizer("../data/finnish.json")
	if err != nil {
		log.Fatal(err)
	}

	// Sample text for embedding
	text := "This is a sample text for embedding. Another sentence for testing."

	// Use SentenceEmbedding function to get embeddings for each sentence
	sentenceEmbeddings, err := embed.SentenceEmbedding(text, sTokenizer, w2v, db, 50)
	if err != nil {
		log.Fatal("Error creating sentence embeddings:", err)
	}

	// Print sentence embeddings
	for i, sentence := range sentenceEmbeddings {
		fmt.Printf("Sentence %d Embedding: %v\n", i+1, sentence.Embedding)
	}

	// Use DocumentEmbedding function to get a single embedding for the entire text
	documentEmbedding, err := embed.DocumentEmbedding(text, w2v, db, 50)
	if err != nil {
		log.Fatal("Error creating document embedding:", err)
	}

	// Print document embedding
	fmt.Println("Document Embedding:", documentEmbedding)
}
