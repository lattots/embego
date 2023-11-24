# Embego

## What is Embego?

Embego is a text embedding library for Go. It has tools to create embeddings for documents and sentences of a document.

## Why Embego?

Embego is the only embedding library for Go that allows you to create sentence and document embeddings instead of just word embeddings. It is written in pure Go, so it is super easy to integrate to Go NLP projects.

## How to install Embego?

For installation, you have 2 options depending on the intended use.

### Install packages as dependencies

```bash
go get github.com/lattots/embego/pkg_name
```

If you install only the packages, you will need to install WeGo to initialize the vector model.

```bash
go get -u github.com/ynqa/wego
bin/wego -h
```

### Install entire application from GitHub

```bash
git clone https://github.com/lattots/embego.git
```

If you install the entire application, you can use ready-made apps like init.go and api.go that make the use easier. If you install the entire application, you can also use packages directly in you own code.

## How to use Embego?

### Initialize vector model

Package parser has tools to initialize new vector model with any text corpus. At the time it only processes data from certain JSON format, but in the future it will handle any kind of text data.

This functionality can be tested with init.go. In order to use init.go you will have to install the entire application from GitHub.

```bash
cd app/create_vectors
go run init.go
```

Now you will have a txt file of tokenized text and an SQLite database with counts of each token in the source data. After this you will need to run script from [WeGo](https://github.com/ynqa/wego):

```bash
wego word2vec -i path/to/txt/file -o path/to/vector/file
```

You can use either Word2Vec, GloVe or LexVec to create the vectors. This script will create a txt file that contains your word embeddings model. Now you are ready to start creating embeddings for words, sentences or documents.

### Create embeddings

Package embed has functionality for creating word, sentence and document embeddings. In order to use this package you must first initialize the vector model. After doing so you should have 3 files in `/data/` directory:

- finnish.json
  - Data for Finnish sentence tokenizer.
  - Data can be downloaded from:
    <https://github.com/neurosnap/sentences/data/>
- word_vectors.txt
  - Embeddings for individual words in your training data.
  - This is created by WeGo
- token_frequency.db
  - Database of all words and their frequencies in the training data.
  - This is created by parser.

Now that you have everything in place, you can start creating word, sentence and document embeddings.

For this you have 2 options. You can use `/app/embeddings_api/api.go` which is a Go API that returns the document and sentence embeddings of any text you send it. The other option is to use embed package directly in your code.

#### Embeddings API

Embeddings API can return requested embeddings as a JSON object.

For example usage, check out `/examples/api_example.go`. In order to run example program, you will have to run `/app/embeddings_api/api.go` at the same time.

#### Embed package

Embed package has 3 exported function:

- Embed
  - This function returns the word embedding of a token you provide it.
- SentenceEmbedding
  - This function returns embeddings for all sentences in the text you provide it.
- DocumentEmbedding
  - This function returns embedding for the entire text you provide it.

For example usage, check out `/examples/embed_example.go`.

## Why bother with embeddings?

### Sentiment analysis and comparisons

With embeddings, it is very easy to compare sentiments (i.e. meanings) of two or more pieces of text. This makes it possible to compare the similarity of two texts without doing just "dumb" word matching.

### Text search

Other use case is text search. This is similar to sentiment analysis but on an even more applied level. With sentence embeddings it is possible to for example search for information from large texts.

Both sentiment analysis and search can be done by calculating the cosine similarity of two or more sentences/documents. Util package has function for calculating just that.
