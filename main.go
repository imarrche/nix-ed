package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func fetchPosts() {
	r, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		log.Fatalf("couldn't fetch posts: %v", err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}

func main() {
	fetchPosts()
}
