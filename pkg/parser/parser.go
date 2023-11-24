package parser

import (
	"fmt"
	"os"

	"github.com/lattots/embego/pkg/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ProcessCorpus(inputFilename string, outputFilename string) (tokens []string, err error) {

	fileContent, err := os.ReadFile(inputFilename)
	if err != nil {
		return nil, err
	}

	corpus := string(fileContent)
	fmt.Println("File read to string")

	// Text corpus is tokenized to a slice of tokens.
	tokens, err = util.Tokenize(corpus)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully tokenized text")

	// Tokens are saved to txt file.
	err = saveTokens(outputFilename, tokens)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Vocabulary saved to file: %v\n", outputFilename)

	// Slice of tokens is returned.
	return tokens, err
}

func InitFrequencyDB(dbFilename string, tokens []string) (TokenCount, error) {
	// Database for keeping track of token frequency is initialized.
	db, err := gorm.Open(sqlite.Open(dbFilename), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Drop "tokens" table if it already exists.
	err = db.Migrator().DropTable(&util.Token{})
	if err != nil {
		return nil, err
	}
	// Schema for database is created.

	err = db.AutoMigrate(&util.Token{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully initialized token count database")

	// Database connection is closed when possible
	defer func() {
		dbInstance, err := db.DB()
		if err != nil {
			panic(err)
		}
		err = dbInstance.Close()
		if err != nil {
			panic(err)
		}
	}()

	tokenCount := make(TokenCount)

	for _, t := range tokens {
		tokenCount[t] += 1
	}
	fmt.Println("Successfully calculated token count")

	err = saveTokenFrequency(db, tokenCount)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully saved token count to database")

	return tokenCount, nil
}
