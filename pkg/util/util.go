// Package util provides utilities for embeddings-api.
package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/neurosnap/sentences"
	"github.com/ynqa/wego/pkg/embedding"
	"gorm.io/gorm"
)

type Token struct {
	Token string `gorm:"primaryKey"`
	Count int
}

type Result struct {
	Tokens []string `json:"lemmatized_words"`
}

type EmptySliceError struct {
	Message string
}

func (e *EmptySliceError) Error() string {
	return e.Message
}

// Tokenize tokenizes text to a slice of words. Removes all special characters from tokens.
func Tokenize(text string) (tokens []string, err error) {

	// All whitespace characters, including newline, are converted to spaces.
	newline := regexp.MustCompile(`\s`)
	text = newline.ReplaceAllString(text, " ")

	// Special characters and numbers are removed from all words.
	special := regexp.MustCompile(`[^\p{L} ]`)
	text = special.ReplaceAllString(text, "")

	// Text is split at spaces.
	words := strings.Fields(strings.ToLower(text))

	// If there are no words to tokenize the function returns an error.
	if len(words) == 0 {
		return nil, &EmptySliceError{"NO WORDS IN TEXT TO TOKENIZE"}
	}

	// Words are lemmatized.
	words, err = lemmatize(words)
	if err != nil {
		return nil, err
	}

	for _, w := range words {
		if w != "" {
			// Cleaned word is added to tokens.
			tokens = append(tokens, w)
		}
	}

	if len(tokens) == 0 {
		return nil, &EmptySliceError{"NO TOKENS IN FINAL TOKEN LIST"}
	}

	// Slice of clean tokens is returned.
	return tokens, nil
}

// GetTokenFrequency gets frequency for given token. Returns frequency and error/nil.
func GetTokenFrequency(token string, db *gorm.DB) (freq int, err error) {
	if err = db.Model(Token{}).Select("count").Where("token = ?", token).Scan(&freq).Error; err != nil {
		return 0, err
	}
	return freq, nil
}

// LoadEmbeddingModel loads embedding model from file. Returns embedding model and error/nil.
func LoadEmbeddingModel(filename string) (embedding.Embeddings, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// Pre-trained embedding model is loaded from file.
	model, err := embedding.Load(file)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// LoadSentenceTokenizer loads sentence tokenizer from file. Returns tokenizer and error/nil.
func LoadSentenceTokenizer(filename string) (*sentences.DefaultSentenceTokenizer, error) {
	// Training data is loaded from file.
	// Use data from https://github.com/neurosnap/sentences/data/
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// Model is trained with the given data.
	training, err := sentences.LoadTraining(data)
	if err != nil {
		return nil, err
	}
	// Tokenizer instance is created with trained model.
	tokenizer := sentences.NewSentenceTokenizer(training)
	return tokenizer, nil
}

// Calculates the length of a vector.
func vectorLength(v []float64) (length float64) {
	var sum float64
	for _, c := range v {
		sum += math.Pow(c, 2)
	}
	return math.Sqrt(sum)
}

// NormalizeVector scales vector a length of 1.
func NormalizeVector(vector []float64) []float64 {
	vecLen := vectorLength(vector)
	for i := range vector {
		vector[i] = vector[i] / vecLen
	}
	return vector
}

// CosineSimilarity calculates the cosine similarity of two vectors.
func CosineSimilarity(v1 []float64, v2 []float64) (similarity float64) {
	var product float64
	for i := range v1 {
		product += v1[i] * v2[i]
	}
	similarity = product / (vectorLength(v1) * vectorLength(v2))

	// Function returns the similarity score of the vectors.
	return similarity
}

// Calls lemmatization API for token lemmatization.
func lemmatize(tokens []string) ([]string, error) {
	// URL of the lemmatization API
	apiURL := "http://localhost:8080/lemmatize"

	// HTTP request payload is json containing slice of tokens.
	payload, err := json.Marshal(map[string]interface{}{"words": tokens})
	if err != nil {
		return nil, err
	}

	// Make a POST request to the API
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Parse the JSON response
	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil, err
	}

	return result.Tokens, nil
}
