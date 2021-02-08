package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
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

	e := echo.New()
	api := e.Group("/api")

	ps := api.Group("/posts")
	ps.GET("", ph.GetAll)
	ps.POST("", ph.Create)
	ps.GET("/:id", ph.GetByID)
	ps.PATCH("/:id", ph.Update)
	ps.DELETE("/:id", ph.DeleteByID)

	cs := api.Group("/comments")
	cs.GET("", ch.GetAll)
	cs.POST("", ch.Create)
	cs.GET("/:id", ch.GetByID)
	cs.PATCH("/:id", ch.Update)
	cs.DELETE("/:id", ch.DeleteByID)

	e.Logger.Fatal(e.Start(":8080"))
}
