package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
	"github.com/shota-tech/graphql/server/graph"
	"github.com/shota-tech/graphql/server/repository"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// connect db
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("failed to load location: %v", err)
	}
	config := mysql.Config{
		DBName:    os.Getenv("MYSQL_DATABASE"),
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Addr:      "db",
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_general_ci",
		Loc:       jst,
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	// DI
	taskRepository := repository.NewTaskRepository(db)
	userRepository := repository.NewUserRepository(db)
	resolver := &graph.Resolver{
		TaskRepository: taskRepository,
		UserRepository: userRepository,
	}

	// start server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", cors.Default().Handler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
