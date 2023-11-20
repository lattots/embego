package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gorm.io/gorm"
)

type Token struct {
	Token string `gorm:"primaryKey"`
	Count int
}

type TokenCount map[string]int

func readCatalog(filename string) (products []Product, err error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// Saves token frequencies to database. Returns error/nil.
func saveTokenFrequency(db *gorm.DB, tokenCount TokenCount) error {
	// Convert token-count pairs to a slice of Token structs
	tokens := make([]Token, 0, len(tokenCount))
	for token, count := range tokenCount {
		tokens = append(tokens, Token{Token: token, Count: count})
	}

	// Use CreateInBatch to perform bulk insert
	r := db.CreateInBatches(tokens, 500)
	if r.Error != nil {
		return r.Error
	}

	// If the save is successful, the function returns nil.
	return nil
}

func formCorpus(products ProductCatalog) string {
	var corpus strings.Builder
	for _, product := range products.Products {
		productText := fmt.Sprintf("%v %v %v", product.Title, product.Desc, product.Gender)
		corpus.WriteString(productText)
		corpus.WriteString(" ")
	}
	return corpus.String()
}

func saveTokens(outputFilename string, tokens []string) error {
	outputTXTFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}

	defer func(outputTXTFile *os.File) {
		err := outputTXTFile.Close()
		if err != nil {

		}
	}(outputTXTFile)

	// Tokens are concatenated to single string
	// Individual tokens are separated by space
	tokensString := strings.Join(tokens, " ")

	// Tokens are written to a text file
	_, err = outputTXTFile.WriteString(tokensString)
	if err != nil {
		return err
	}

	// If write is successful, the function returns nil
	return nil
}
