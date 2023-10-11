# Embego

## What is embego?

Embego is a text embedding library for Go. It has tools to create embeddings for documents and sentences of a document.

## Why embego?

Embego is the only embedding library for Go that allows you to create sentence and document embeddings instead of just word embeddings. It is written in pure Go so it is super easy to integrate to Go NLP projects.

## How to install embego?

For installation you have 2 options depending on the intended use.

### Install packages as dependencies

```bash
go get github.com/lattots/embego/pkg_name
```

If you install only the packages, you will need to install wego to initialize the vector model.

```bash
go get -u github.com/ynqa/wego
bin/wego -h
```

### Install entire application from github

```bash
git clone https://github.com/lattots/embego.git
```

If you install the entire application, you can use ready made apps like init.go and api.go that make the use easier. If you install the entire application, you can also use packages directly in you own code.

## How to use embego?

### Initialize vector model

Package parser has tools to initialize new vector model with any text corpus. At the time it only processes data from certain JSON format, but in the future it will handle any kind of text data.

This functionality can be tested with init.go. In order to use init.go you will have to install the entire application from github.

```bash
cd app/create_vectors
go run init.go
```

Now you will have a txt file of tokenized text and an SQLite database with counts of each token in the source data. After this you will need to run script from [wego](https://github.com/ynqa/wego):

```bash
wego word2vec -i path/to/txt/file -o path/to/vector/file
```

You can use either Word2Vec, GloVe or LexVec to create the vectors. This script will create a txt file that contains your word embeddings model. Now you are ready to start creating embeddings for words, sentences or documents.

### Create embeddings

Package embed has functionality for creating word, sentence and document embeddings. In order to use this package you must first initialize vector model. After doing so you should have 3 files in /data/ directory:

- finnish.json
  - Data for Finnish sentence tokenizer.
  - Data can be downloaded from:
    <https://github.com/neurosnap/sentences/data/>
- word_vectors.txt
  - Embeddings for individual words in you training data.
  - This is created by wego
- token_frequency.db
  - Database of all words and their frequencies in the training data.
  - This is created by parser.

Now that you have everything in place, you can start creating word, sentence and document embeddings.

For this you have 2 options. You can use /app/embeddings_api/api.go which is a Go API that returns the document and sentence embeddings of any text you send it. The other option is to use embed package directly in your code.

#### Embeddings API

Example usage:

```go
type Response struct {
  DocumentEmbedding  []float64   `json:"document_embedding"`
  SentenceEmbeddings [][]float64 `json:"sentence_embeddings"`
}

var embeddings []Response

texts := []string{
    "Tämä merkkijono sisältää esimerkkitekstiä. Tätä voidaan käyttää vektorirajapinnan testaamiseen.",
    "Tämä on toinen merkkijono. Vektorirajapinta palauttaa jokaista merkkijonoa vastaavan dokumenttivektorin, mutta myös jokaista virkettä vastaavan virkevektorin.",
    "Yhteensä rajapinta palauttaa siis yhden dokumenttivektorin ja listan virkevektoreita.",
}

for _, text := range texts {
    body, err := json.Marshal(map[string]string{
        "text": text,
    })
    if err != nil {
        panic(err)
    }
    resp, err := http.Post("http://localhost:8081/embedding", "application/json", bytes.NewReader(body))
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    jsonString, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    var response Response
    err = json.Unmarshal(jsonString, &response)
    if err != nil {
        panic(err)
    }

    embeddings = append(embeddings, response)
}
```

#### Embed package

Embed package has 3 exported function:

- Embed
  - This function returns the word embedding of a token you provide it.
- SentenceEmbedding
  - This function returns embeddings for all sentences in the text you provide it.
- DocumentEmbedding
  - This function returns embedding for the entire text you provide it.

Example usage:

```go

import (
  "github.com/lattots/embego/pkg/embed"
  "github.com/lattots/embego/pkg/util"
  "gorm.io/driver/sqlite"
  "gorm.io/gorm"
)

// Load embedding model from word vectors file.
eModel, err := util.LoadEmbeddingModel("/paht/to/vector/model")
if err != nil {
    panic(err)
}

// Load sentence tokenizer if you want to create sentence embeddings.
sTokenizer, err := util.LoadSentenceTokenizer("/path/to/sentence/tokenizer/model")
if err != nil {
    panic(err)
}

// Load token frequency database if you want to create sentence or document embeddings.
db, err := gorm.Open(sqlite.Open("/path/to/token/frequency/database"), &gorm.Config{})
if err != nil {
    panic(err)
}

// Declare constant "a". This is a parametre for embedding algorithm. "a" affects the significance of token frequency in creatin multi word embeddings. If "a" is small, the significance of the number of appearances of a token in training data is large. If "a" is large, the significance is greater.
const a float64 = 3

// This is how you create embedding for entire text.
docEmbedding, err := embed.DocumentEmbedding(text, eModel, db, a)
if err != nil {
    panic(err)
}

// This is how you create embeddings for all sentences in text.
sentences, err := embed.SentenceEmbedding(text, sTokenizer, eModel, db, a)
if err != nil {
    panic(err)
}

// Note that sentences is a slice of type Sentence. See /pkg/embed/main.go for details.
// You can access individual embeddings like this.
singleEmbedding := sentences[0].Embedding
```

## What can you do with embeddings?

### Sentiment analysis and comparisons

With embeddings it is very easy to compare sentiments (ie. meanings) of two or more pieces of text. This makes it possible to compare the similarity of two texts without doing just dumb "word matching".

### Text search

Other use case is text search. This is similar to sentiment analysis but on an even more applied level. With sentence embeddings it is possible to for example search for information from large texts.

Both sentiment analysis and search can be done by calculating the cosine similarity of two or more sentences/documents. Util package has function for calculating just that.

```go
import "github.com/lattots/embego/pkg/util"

# Create embeddings for two or more sentences/documents
var v1, v2 []float64

similarity := util.CosineSimilarity(v1, v2)
```
