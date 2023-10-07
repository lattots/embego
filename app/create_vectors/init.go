package main

import (
	"github.com/lattots/embego/pkg/parser"
)

const (
	inputFilename    = "/home/otso/go/src/embeddings-api/data/varuste-net_products.json"
	outputFilename   = "/home/otso/go/src/embeddings-api/data/vocabulary.txt"
	databaseFilename = "/home/otso/go/src/embeddings-api/data/token_frequency.db"
)

func main() {
	// Tokens are processed from input data.
	tokens, err := parser.ProcessCorpus(inputFilename, outputFilename)
	if err != nil {
		panic(err)
	}

	// Token frequencies are saved to SQLite database.
	_, err = parser.InitFrequencyDB(databaseFilename, tokens)
	if err != nil {
		panic(err)
	}
}
