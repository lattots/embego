// Running this program will start a local embeddings API that other processes can request text embeddings from.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/neurosnap/sentences"
	"github.com/ynqa/wego/pkg/embedding"
	"io"
	"log"
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

	fmt.Println("Running api.")

	// Word2Vec embeddings model is loaded.
	w2v, err := util.LoadEmbeddingModel("../../data/english/w2v.txt")
	if err != nil {
		log.Fatal(err)
	}

	// GloVe embeddings model is loaded.
	glove, err := util.LoadEmbeddingModel("../../data/english/glove.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Token frequency database connection is initialized.
	db, err := gorm.Open(sqlite.Open("../../data/english/token_frequency.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Sentence tokenizer is loaded.
	st, err := util.LoadSentenceTokenizer("../../data/english.json")
	if err != nil {
		log.Fatal(err)
	}

	// Handles requests to the Word2Vec endpoint.
	http.HandleFunc("/embedding/word2vec", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, w2v, db, st)
	})

	// Handles requests to the GloVe endpoint.
	http.HandleFunc("/embedding/glove", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, glove, db, st)
	})

	fmt.Printf("Server started on port :%v\n", PORT)
	err = http.ListenAndServe(":"+fmt.Sprint(PORT), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Handles embeddings request with specified model.
func handleRequest(w http.ResponseWriter, r *http.Request, eModel embedding.Embeddings, db *gorm.DB, st *sentences.DefaultSentenceTokenizer) {

	// Request body is read.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request")
		log.Fatal(err)
	}

	// Request JSON is converted to a map called payload.
	var payload map[string]string
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("Error unmarshalling request body")
		log.Fatal(err)
	}

	// Text field from payload is read.
	text := payload["text"]

	if len(text) == 0 {
		http.Error(w, "Text field is required", http.StatusBadRequest)
		return
	}

	// This is the response that is returned to client.
	var response Response

	// Embedding for entire text is created.
	docEmbedding, err := embed.DocumentEmbedding(text, eModel, db, a)
	if err != nil {
		log.Println("Error creating document embedding")
		log.Fatal(err)
	}
	// Document embedding is added to response.
	response.DocumentEmbedding = docEmbedding

	// Embeddings for all sentences are created.
	senEmbeddings, err := embed.SentenceEmbedding(text, st, eModel, db, a)
	if err != nil {
		log.Println("Error creating sentence embeddings")
		log.Fatal(err)
	}
	// Sentence embeddings are added to response.
	for _, s := range senEmbeddings {
		response.SentenceEmbeddings = append(response.SentenceEmbeddings, s.Embedding)
	}

	// Write the JSON response to the client.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {

		log.Fatal(err)
	}
}
