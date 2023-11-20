package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/lattots/embego/pkg/util"
)

type Response struct {
	DocumentEmbedding  []float64   `json:"document_embedding"`
	SentenceEmbeddings [][]float64 `json:"sentence_embeddings"`
}

const prodsPath string = "../data/product_info.csv"

func main() {
	products, err := util.ReadProductTable(prodsPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully read product table")

	for i, p := range products {
		fmt.Println("Processing: ", i)
		p.Word2Vec, err = embed(p.Desc, "http://localhost:8081/embedding/word2vec")
		if err != nil {
			log.Fatal(err)
		}

		p.GloVe, err = embed(p.Desc, "http://localhost:8081/embedding/glove")
		if err != nil {
			log.Fatal(err)
		}

		products[i] = p
	}

	err = util.UpdateProductTable(prodsPath, products)
	if err != nil {
		log.Fatal(err)
	}
}

func embed(text string, url string) ([]float64, error) {
	body, err := json.Marshal(map[string]string{
		"text": text,
	})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		// Print the error message from the response body
		errMsg, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		log.Println("Error: ", string(errMsg))
		return nil, err
	}

	jsonString, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(jsonString, &response)
	if err != nil {
		return nil, err
	}

	return response.DocumentEmbedding, nil
}
