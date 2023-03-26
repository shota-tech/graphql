package main

//go:generate gqlgen generate
//go:generate sqlboiler mysql

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-sql-driver/mysql"
	"github.com/shota-tech/graphql/server/graph"
	"github.com/shota-tech/graphql/server/loader"
	"github.com/shota-tech/graphql/server/middleware/auth"
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
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	todoRepository := repository.NewTodoRepository(db)
	taskLoader := loader.NewTaskLoader(taskRepository)
	userLoader := loader.NewUserLoader(userRepository)
	todoLoader := loader.NewTodoLoader(todoRepository)
	loaders := loader.NewLoaders(
		userLoader,
		taskLoader,
		todoLoader,
	)
	resolver := &graph.Resolver{
		Loaders:        loaders,
		UserRepository: userRepository,
		TaskRepository: taskRepository,
		TodoRepository: todoRepository,
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// setup router
	router := chi.NewRouter()
	router.Use(chiMiddleware.AllowContentType("application/json"))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))
	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.With(auth.EnsureValidToken()).Handle("/graphql", srv)

	// start server
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
