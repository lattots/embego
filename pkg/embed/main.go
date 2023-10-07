// Provides functions for creating embeddings.
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
	Length     int
	Embeddings []Embedding
	Embedding  []float64
}

type Embedding struct {
	Token  string
	Vector []float64
}

type vec = vector.Vector

// Creates embeddings for each sentence in text. Returns a slice of type Sentence.
func SentenceEmbedding(text string, sTokenizer *sentences.DefaultSentenceTokenizer, eModel embedding.Embeddings, db *gorm.DB, a float64) ([]Sentence, error) {
	// Tokenize the input text into sentences
	sentencesText := sTokenizer.Tokenize(text)

	// Slice of individual sentences from text
	var sentences []Sentence

	// Iterate over all sentences
	for _, sentenceText := range sentencesText {
		// Sentence is split into a slice of tokens
		tokens, err := util.Tokenize(sentenceText.Text)
		if err != nil {
			return nil, err
		}

		var sentence Sentence

		// Iterate over all tokens
		for _, token := range tokens {
			// Embedding is retrieved for token
			embedding, err := Embed(eModel, token)
			if err != nil {
				log.Printf("Embedding for %v doesn't exist.\n", token)
			} else {
				sentence.Embeddings = append(sentence.Embeddings, embedding)
			}
		}
		sentence.Length = len(tokens)
		// Sentence is appended to the slice of sentences
		sentences = append(sentences, sentence)
	}

	for i := range sentences {
		embedding, err := multiWordVec(sentences[i].Embeddings, db, a)
		if err != nil {
			return nil, err
		}
		sentences[i].Embedding = util.NormalizeVector(embedding)
	}
	return sentences, nil
}

func DocumentEmbedding(text string, eModel embedding.Embeddings, db *gorm.DB, a float64) ([]float64, error) {
	// Text is split to a slice of tokens.
	tokens, err := util.Tokenize(text)
	if err != nil {
		return nil, err
	}

	// Slice that contains all word embeddings of text. This will be populated later.
	var wordEmbeddings []Embedding

	for _, token := range tokens {
		// Embedding is retrieved for token.
		embedding, err := Embed(eModel, token)
		if err != nil {
			log.Printf("Embedding for %v doesn't exist.\n", token)
		} else {
			// Created embedding is appended to the slice.
			wordEmbeddings = append(wordEmbeddings, embedding)
		}
	}

	// Embedding for text is created.
	embedding, err := multiWordVec(wordEmbeddings, db, a)
	if err != nil {
		return nil, err
	}
	embedding = util.NormalizeVector(embedding)

	return embedding, nil
}

// Creates single embedding from all embeddings. Returns single embedding and error/nil.
func multiWordVec(embeddings []Embedding, db *gorm.DB, a float64) (vector []float64, err error) {
	vector = make([]float64, len(embeddings[0].Vector))
	for _, embedding := range embeddings {
		tokenFreq, err := util.GetTokenFrequency(embedding.Token, db)
		if err != nil {
			return nil, err
		}
		scaler := a / (a + float64(tokenFreq))
		tokenVector := vec(embedding.Vector)
		vector = vec.Sum(vector, tokenVector.Scale(scaler))
	}
	vector = vec.Scale(vector, float64(1)/float64(len(embeddings)))
	return vector, nil
}

// Creates single embedding for token. Returns Embedding and error/nil.
func Embed(eModel embedding.Embeddings, token string) (Embedding, error) {
	// Vector corresponding to token is retrieved from eModel.
	vec, exists := eModel.Find(token)

	if exists {
		// Embedding instance is created with token and vector.
		embedding := Embedding{
			Token:  vec.Word,
			Vector: vec.Vector,
		}
		return embedding, nil
	} else {
		// If embedding for token doesn't exist, the function returns non fatal error.
		return Embedding{}, fmt.Errorf("Embedding for %v doesn't exist", token)
	}
}
