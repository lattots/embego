package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Response struct {
	DocumentEmbedding  []float64   `json:"document_embedding"`
	SentenceEmbeddings [][]float64 `json:"sentence_embeddings"`
}

func main() {
	// Define the API endpoint URLs
	word2vecURL := "http://localhost:8081/embedding/word2vec"
	gloveURL := "http://localhost:8081/embedding/glove"

	// Sample text for embedding
	text := "I am trying to learn program development. I am super interested in artificial intelligence."

	// Call Word2Vec API
	word2vecResponse, err := callAPI(word2vecURL, text)
	if err != nil {
		log.Fatal("Error calling Word2Vec API:", err)
	}
	fmt.Println("Word2Vec Response:", word2vecResponse)

	// Call Glove API
	gloveResponse, err := callAPI(gloveURL, text)
	if err != nil {
		log.Fatal("Error calling Glove API:", err)
	}
	fmt.Println("Glove Response:", gloveResponse)
}

func callAPI(apiURL string, text string) (Response, error) {
	// Create a map for the request payload
	payload := map[string]string{"text": text}

	// Convert payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return Response{}, err
	}

	// Make a POST request to the API
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("API returned non-OK status: %v", resp.Status)
	}

	// Unmarshal the JSON response
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}
