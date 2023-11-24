# Example usage

## API example

`api_example.go` is meant for trying out `embeddings_api.go`. Running the program requires you to also run the API program at the same time. The example program embeds one sentence with both Word2Vec and GloVe endpoints. Before you can run the API you must have initialized Word2Vec and GloVe vectors for your vocabulary. See root level README for instructions.

Run the API with:

~~~bash
go run api.go
~~~

Run the example program with:

~~~bash
go run embeddings_api.go
~~~

## Embed example

`embed_example.go` is meant for trying out the embed package of Embego. Running the program creates document and sentence embeddings for a short text sample. Before running the program you must have initialized Word2Vec vectors for your vocabulary. See root level README for instructions.

~~~bash
go run embed_example.go
~~~
