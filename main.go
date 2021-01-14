package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/imarrche/nix-ed/pkg/model"
)

var (
	createPostsTableQuery = `
		CREATE TABLE posts(
			id int primary key,
			title text,
			body text,
			user_id int
		);
	`
	createCommentsTableQuery = `
		CREATE TABLE comments(
			id int primary key,
			name text,
			email text,
			body text,
			post_id int
		);
	`
	insertPostQuery = `
		INSERT INTO posts (id, title, body, user_id) VALUES (?, ?, ?, ?);
	`
	insertCommentQuery = `
		INSERT INTO comments (id, name, email, body, post_id) VALUES (?, ?, ?, ?, ?);
	`
)

func migrate(db *sql.DB) error {
	if _, err := db.Exec("DROP TABLE IF EXISTS posts;"); err != nil {
		return err
	}
	if _, err := db.Exec(createPostsTableQuery); err != nil {
		return err
	}
	if _, err := db.Exec("DROP TABLE IF EXISTS comments;"); err != nil {
		return err
	}
	if _, err := db.Exec(createCommentsTableQuery); err != nil {
		return err
	}
	return nil
}

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

func savePost(db *sql.DB, p model.Post) error {
	if _, err := db.Exec(insertPostQuery, p.ID, p.Title, p.Body, p.UserID); err != nil {
		return err
	}
	return nil
}

func saveComments(db *sql.DB, ch chan model.Comment, done chan struct{}) {
	for c := range ch {
		if _, err := db.Exec(insertCommentQuery, c.ID, c.Name, c.Email, c.Body, c.PostID); err != nil {
			log.Fatal(err)
		}
	}

	done <- struct{}{}
}

func main() {
	db, err := sql.Open("mysql", "root:password@/posts")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := migrate(db); err != nil {
		log.Fatal(err)
	}

	ch := make(chan model.Comment)
	done := make(chan struct{})
	posts, err := fetchPostsByUserID(7)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range posts {
		savePost(db, p)
	}

	go fetchComments(posts, ch)
	go saveComments(db, ch, done)

	<-done
}
