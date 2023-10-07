package parser

import (
	"fmt"

	"github.com/lattots/embego/pkg/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	Id     int     `json:"id"`
	Title  string  `json:"title"`
	Desc   string  `json:"desc"`
	Price  float64 `json:"price"`
	Gender string  `json:"gender"`
}

type ProductCatalog struct {
	Products []Product
}

func ProcessCorpus(inputFilename string, outputFilename string) (tokens []string, err error) {
	// Product catalog is read from json file
	products, err := readCatalog(inputFilename)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully read product catalog")
	catalog := ProductCatalog{
		Products: products,
	}

	// Text corpus is formed from product information.
	corpus := formCorpus(catalog)

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
	db.Migrator().DropTable(&util.Token{})
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

	var tokenCount TokenCount = make(TokenCount)

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
