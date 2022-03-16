package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	h "github.com/bhanupbalusu/ecommerce_v0.1/HEX-MICROSERVICE/api"
	mr "github.com/bhanupbalusu/ecommerce_v0.1/HEX-MICROSERVICE/repository/mongodb"
	shortener "github.com/bhanupbalusu/ecommerce_v0.1/HEX-MICROSERVICE/urlshortener"
)

func main() {
	repo := getRepo()
	fmt.Println(repo)
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func getRepo() shortener.RedirectRepository {
	fmt.Println("............Getting mongodb repository .............")
	repo, err := mr.NewMongoRepository("mongodb://localhost:27017", "shortener", 30)
	if err != nil {
		log.Fatal(err)
	}
	return repo

}
