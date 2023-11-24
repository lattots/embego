// Package embed provides functions for creating embeddings.
package embed

import (
	"fmt"
	"log"

	"github.com/lattots/embego/pkg/util"
	"github.com/neurosnap/sentences"
	"github.com/quartercastle/vector"
	"github.com/ynqa/wego/pkg/embedding"
	"gorm.io/gorm"
)

type Sentence struct {
	Tokens     []string
	Length     int
	Embeddings []Embedding
	Embedding  []float64
}

type Embedding struct {
	Token  string
	Vector []float64
}

type vec = vector.Vector

// SentenceEmbedding creates embeddings for each sentence in text. Returns a slice of type Sentence and error.
func SentenceEmbedding(text string, sTokenizer *sentences.DefaultSentenceTokenizer, eModel embedding.Embeddings, db *gorm.DB, a float64) ([]Sentence, error) {

	// Tokenize the input text into a slice of sentences
	textSens := sTokenizer.Tokenize(text)

	// Slice of individual sentences from text
	var sens []Sentence

	// Each sentence is tokenized and added to the slice of sentences.
	for _, sen := range textSens {
		tokens, err := util.Tokenize(sen.Text)
		if err != nil {
			log.Println("Error tokenizing sentence: ", sen)
			return nil, err
		}

		var s Sentence
		s.Tokens = tokens
		s.Length = len(tokens)

		// Sentences with 0 tokens are skipped. They contain no information and will cause problems later on.
		if s.Length == 0 {
			continue
		}

		sens = append(sens, s)
	}

	// Embeddings for all tokens in all sentences are created.
	for i := range sens {
		sens[i].createEmbeddings(eModel)
	}

	for i := range sens {
		if len(sens[i].Embeddings) == 0 {
			log.Println("No embeddings for sentence: ", sens[i])
			continue
		}
		e, err := multiWordVec(sens[i].Embeddings, db, a)
		if err != nil {
			return nil, err
		}
		sens[i].Embedding = util.NormalizeVector(e)
	}

	return sens, nil
}

// DocumentEmbedding creates single embedding for the text provided. Returns slice of float64 and error.
func DocumentEmbedding(text string, eModel embedding.Embeddings, db *gorm.DB, a float64) ([]float64, error) {

	// Text is split to a slice of tokens.
	tokens, err := util.Tokenize(text)
	if err != nil {
		return nil, err
	}

	// Slice that contains all word embeddings of text. This will be populated later.
	var wordEmbeddings []Embedding

	// Embedding for each token is created and added to the slice of embeddings.
	for _, token := range tokens {
		// Embedding is retrieved for token.
		wordE, err := Embed(eModel, token)
		if err != nil {
			log.Printf("Embedding for %v doesn't exist.\n", token)
		} else {
			// Created embedding is appended to the slice.
			wordEmbeddings = append(wordEmbeddings, wordE)
		}
	}

	// Embedding for text is created.
	docE, err := multiWordVec(wordEmbeddings, db, a)
	if err != nil {
		log.Printf("Error creating multi word vector")
		return nil, err
	}
	docE = util.NormalizeVector(docE)

	return docE, nil
}

// Embed creates single embedding for token. Returns Embedding and error.
func Embed(eModel embedding.Embeddings, token string) (Embedding, error) {

	// Vector corresponding to token is retrieved from eModel.
	vec, exists := eModel.Find(token)

	if exists {
		// Embedding instance is created with token and vector.
		e := Embedding{
			Token:  vec.Word,
			Vector: vec.Vector,
		}
		return e, nil
	} else {
		// If embedding for token doesn't exist, the function returns non-fatal error.
		return Embedding{}, fmt.Errorf("EMBEDDING FOR %v DOESN'T EXIST", token)
	}
}

// Creates embeddings for all tokens in a sentence.
func (s *Sentence) createEmbeddings(eModel embedding.Embeddings) {

	// Embedding for each token in a sentence is created and added to the embeddings of sentence object.
	for _, t := range s.Tokens {
		e, err := Embed(eModel, t)
		if err != nil {
			log.Printf("Embedding for %v doesn't exist.\n", t)
			continue
		}

		s.Embeddings = append(s.Embeddings, e)
	}
}

// Creates single embedding from all embeddings. Returns single embedding and error.
func multiWordVec(embeddings []Embedding, db *gorm.DB, a float64) (vector []float64, err error) {

	// Final vector that will contain the multi-word vector.
	vector = make([]float64, len(embeddings[0].Vector))

	if len(embeddings) == 0 {
		log.Fatal("Error creating multi word vector for slice of length 0")
	}

	// Multi-word vector for all embeddings is calculated.
	for _, e := range embeddings {
		tokenFreq, err := util.GetTokenFrequency(e.Token, db)
		if err != nil {
			return nil, err
		}
		scalar := a / (a + float64(tokenFreq))
		tokenVector := vec(e.Vector)
		vector = vec.Sum(vector, tokenVector.Scale(scalar))
	}
	vector = vec.Scale(vector, float64(1)/float64(len(embeddings)))

	return vector, nil
}
