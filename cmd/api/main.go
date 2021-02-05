package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/imarrche/nix-ed/internal/comment"
	"github.com/imarrche/nix-ed/internal/model"
	"github.com/imarrche/nix-ed/internal/post"
)

func main() {
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&model.Post{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&model.Comment{}); err != nil {
		log.Fatal(err)
	}

	ph := post.NewHandler(post.NewService(post.NewRepo(db)))
	ch := comment.NewHandler(comment.NewService(comment.NewRepo(db)))
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/posts", func(r chi.Router) {
			ph.Init(r)
		})
		r.Route("/comments", func(r chi.Router) {
			ch.Init(r)
		})
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("sarting the server")
	s.ListenAndServe()
}
