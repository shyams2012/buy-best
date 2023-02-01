package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/shyams2012/buy-best/graph/generated"
	"github.com/shyams2012/buy-best/graph/interfaces"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultPort = "8800"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("[WARNING] %v", err)
		log.Printf("[WARNING] %v", ".env file doesn't exist, so please provide environment variables from some other way using .env.example as a reference ")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = defaultPort
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DB"))
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	resolver := interfaces.NewResolver(dsn, config, true)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	mux := http.NewServeMux()

	mux.Handle("/__playground", playground.Handler("GraphQL playground", "/query"))

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/query", interfaces.MiddlewareGetUserFromToken(srv, resolver))

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PATCH", "PUT", "DELETE", "HEAD"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Requested-With", "Content-Type", "Authorization", "Accept-Encoding", "Host", "Origin", "Accept"},
		Debug:            false,
	})
	muxHandler := c.Handler(mux)

	log.Printf("connect to http://localhost:%s/__playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, muxHandler))

}
