// Running this program will initialize new embedding model training data to the data directory.
package main

import (
	"github.com/lattots/embego/pkg/parser"
)

const (
	inputFilename    = "../../data/english/corpus.txt"
	outputFilename   = "../../data/english/vocabulary.txt"
	databaseFilename = "../../data/english/token_frequency.db"
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
