package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/imarrche/nix-ed/pkg/model"
)

func fetchPostsByUserID(id int) (ps []model.Post, err error) {
	r, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts?userId=%d", id))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&ps)
	return
}

func fetchComments(posts []model.Post, ch chan model.Comment) {
	var cs []model.Comment

	for _, p := range posts {
		r, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/comments?postId=%d", p.ID))
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()

		json.NewDecoder(r.Body).Decode(&cs)
		for _, c := range cs {
			ch <- c
		}
	}

	close(ch)
}

func main() {
	db, err := gorm.Open(mysql.Open("root:password@/posts"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.Post{})
	db.AutoMigrate(&model.Comment{})

	ch := make(chan model.Comment)
	done := make(chan struct{})
	posts, err := fetchPostsByUserID(7)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range posts {
		db.Create(&p)
	}

	go fetchComments(posts, ch)
	go func() {
		for c := range ch {
			db.Create(&c)
		}

		done <- struct{}{}
	}()

	<-done
}
