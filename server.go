package main

import (
	"context"
	"errors"
	"log"
	"myapp/config"
	"myapp/directives"
	"myapp/graph"
	"myapp/middleware"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"

	"myapp/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	config.ConnectDB()
	db := config.GetDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	router := chi.NewRouter()
	router.Use(
		middleware.AuthMiddleware,
		middleware.CORSMiddleware,
	)

	c := generated.Config{Resolvers: &graph.Resolver{}}
	c.Directives.IsLogin = directives.IsLogin

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.Printf("PANIC ERROR: %+v\n", err)
		debug.PrintStack()

		switch reflect.TypeOf(err).String() {
		case "*gqlerror.Error":
			return err.(error)
		}

		return errors.New("internal system error")
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
