package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ahsan-arif/go-rss-scrapper/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}
	fmt.Println("port: ", portString)

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("dbURL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Connection to database failed ", dbURL)
	}
	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}
	go startScrapping(db, 10, time.Minute)

	router := chi.NewRouter()

	//setting up cors

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//create a new router to handle requests
	v1Router := chi.NewRouter()
	v1Router.Get("/ready", handleRequest)
	v1Router.Get("/error", handleErr)
	v1Router.Post("/users/", apiCfg.handlerCreateUser)
	v1Router.Get("/users/", apiCfg.middlewreAuth(apiCfg.handlerGetUser))

	v1Router.Post("/feeds/", apiCfg.middlewreAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds/", apiCfg.handlerGetFeeds)

	v1Router.Get("/posts", apiCfg.middlewreAuth(apiCfg.handlerGetPostsForUser))

	v1Router.Post("/feed_follows/", apiCfg.middlewreAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows/", apiCfg.middlewreAuth(apiCfg.handlerGetFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewreAuth(apiCfg.handlerDeleteFeedFollow))

	//nest the router within the existing router
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
