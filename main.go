package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var n = 100

func fetchPostByID(id int, ch chan []byte) {
	r, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id))
	if err != nil {
		log.Fatalf("couldn't fetch post %d: %v", id, err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	ch <- body
}

func main() {
	ch := make(chan []byte)
	for id := 1; id <= n; id++ {
		go fetchPostByID(id, ch)
	}

	for i := 1; i <= n; i++ {
		// ioutil
		err := ioutil.WriteFile(fmt.Sprintf("./storage/posts/%d.txt", i), <-ch, 0755)
		if err != nil {
			log.Fatalf("couldn't write file %d: %v", i, err)
		}

		// bufio
		// Can write into anything that implements io.Writer interface.
		// f, err := os.Create(fmt.Sprintf("./storage/posts/%d.txt", i))
		// if err != nil {
		// 	log.Fatalf("couldn't create file %d: %v", i, err)
		// }
		// f.Sync()
		// w := bufio.NewWriter(f)
		// if _, err := w.Write(<-ch); err != nil {
		// 	log.Fatalf("couldn't write file %d: %v", i, err)
		// }
		// w.Flush()
	}
}
