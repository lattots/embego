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
	defer jsonFile.Close()

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
	// All token counts are saved to the database provided.
	for token, count := range tokenCount {
		r := db.Create(&Token{Token: token, Count: count})
		if r.Error != nil {
			return r.Error
		}
	}
	// If the save is successful, the funtion returns nil.
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
	// Tokens are concatenated to single string
	// Individual tokens are separated by space
	tokensString := strings.Join(tokens, " ")

	// Tokens are written to a text file
	err := os.WriteFile(outputFilename, []byte(tokensString), 0644)
	if err != nil {
		return err
	}
	// If the write is successful, the function returns nil
	return nil
}
