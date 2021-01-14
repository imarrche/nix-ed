package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var n = 100

func fetchPostByID(id int, ch chan string) {
	r, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id))
	if err != nil {
		log.Fatalf("couldn't fetch post %d: %v", id, err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	ch <- string(body)
}

func main() {
	ch := make(chan string)
	for id := 1; id <= n; id++ {
		go fetchPostByID(id, ch)
	}

	for i := 1; i <= n; i++ {
		log.Println(<-ch)
	}
}
