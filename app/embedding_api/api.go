package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lattots/embego/pkg/embed"
	"github.com/lattots/embego/pkg/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const PORT int = 8081
const a float64 = 50

type Response struct {
	DocumentEmbedding  []float64   `json:"document_embedding"`
	SentenceEmbeddings [][]float64 `json:"sentence_embeddings"`
}

func main() {
	eModel, err := util.LoadEmbeddingModel("/home/otso/go/src/embeddings-api/data/word_vectors.txt")
	if err != nil {
		panic(err)
	}

	sTokenizer, err := util.LoadSentenceTokenizer("/home/otso/go/src/embeddings-api/data/finnish.json")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("/home/otso/go/src/embeddings-api/data/token_frequency.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/embedding", func(w http.ResponseWriter, r *http.Request) {
		// Request body is read.
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		// Request JSON is converted to a map called payload.
		var payload map[string]string
		err = json.Unmarshal(body, &payload)
		if err != nil {
			panic(err)
		}

		// Text field from payload is read.
		text := payload["text"]

		// This is the response that is returned to client.
		var response Response

		// Embedding for entire document is created.
		docEmbedding, err := embed.DocumentEmbedding(text, eModel, db, a)
		if err != nil {
			panic(err)
		}
		// Embedding is added to response.
		response.DocumentEmbedding = docEmbedding

		// Embedding for each sentence is the text is created.
		sentences, err := embed.SentenceEmbedding(text, sTokenizer, eModel, db, a)
		if err != nil {
			panic(err)
		}

		// Embedding of each sentence is added to response.
		for _, s := range sentences {
			response.SentenceEmbeddings = append(response.SentenceEmbeddings, s.Embedding)
		}

		// Marshal the response to JSON.
		// jsonResponse, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}

		// Write the JSON response to the client.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	fmt.Printf("Server started on port :%v", PORT)
	http.ListenAndServe(":"+fmt.Sprint(PORT), nil)
}
