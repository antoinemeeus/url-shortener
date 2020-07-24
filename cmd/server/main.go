package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/antoinemeeus/url-shortener/pkg/api"
	"github.com/antoinemeeus/url-shortener/pkg/shortener"
	sr "github.com/antoinemeeus/url-shortener/pkg/storage/redis"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := api.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)
	r.Put("/", handler.Update)

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

func chooseRepo() shortener.RedirectRepository {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	switch os.Getenv("DB_ENGINE") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := sr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
		// case "mongo":
		// 	mongoURL := os.Getenv("MONGO_URL")
		// 	mongodb := os.Getenv("MONGO_DB")
		// 	mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		// 	repo, err := mr.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	return repo
	}
	return nil
}
