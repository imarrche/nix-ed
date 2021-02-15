package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/imarrche/nix-ed/docs"
	"github.com/imarrche/nix-ed/internal/auth"
	"github.com/imarrche/nix-ed/internal/comment"
	"github.com/imarrche/nix-ed/internal/model"
	"github.com/imarrche/nix-ed/internal/post"
)

// @title Nix-Ed REST API
// @version 1.0
// @description This is simple REST API with CRUD for posts and comments.
// @host localhost:8080
// @BasePath /api/
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

	as := auth.NewGoogleService()
	ph := post.NewHandler(post.NewService(post.NewRepo(db)), as)
	ch := comment.NewHandler(comment.NewService(comment.NewRepo(db)), as)
	ah := auth.NewHandler(as)

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := e.Group("/auth")
	auth.GET("/google/sign-in", ah.GoogleSignIn)
	auth.GET("/google/callback", ah.GoogleCallback)

	api := e.Group("/api")

	ps := api.Group("/posts")
	ps.GET("", ph.GetAll)
	ps.POST("", ph.Create, ph.Auth)
	ps.GET("/:id", ph.GetByID)
	ps.PATCH("/:id", ph.Update, ph.Auth, ph.PostAuthor)
	ps.DELETE("/:id", ph.DeleteByID, ph.Auth, ph.PostAuthor)

	cs := api.Group("/comments")
	cs.GET("", ch.GetAll)
	cs.POST("", ch.Create, ch.Auth)
	cs.GET("/:id", ch.GetByID)
	cs.PATCH("/:id", ch.Update, ch.Auth, ch.CommentAuthor)
	cs.DELETE("/:id", ch.DeleteByID, ch.Auth, ch.CommentAuthor)

	e.Logger.Fatal(e.Start(":8080"))
}
