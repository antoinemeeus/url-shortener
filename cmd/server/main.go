package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/antoinemeeus/url-shortener/pkg/api"
	"github.com/antoinemeeus/url-shortener/pkg/shortener"
	ms "github.com/antoinemeeus/url-shortener/pkg/storage/mysql"
	ps "github.com/antoinemeeus/url-shortener/pkg/storage/postgresql"
	sr "github.com/antoinemeeus/url-shortener/pkg/storage/redis"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	repo := chooseRepo()
	defer repo.Close()
	service := shortener.NewRedirectService(repo)
	handler := api.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/create", handler.Post)
	r.Put("/update", handler.Update)

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

	fmt.Printf("Terminated by %s", <-errs)
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
	timeout, _ := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
	switch os.Getenv("DB_ENGINE") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := sr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "postgres":
		psqlHost := os.Getenv("POSTGRESQL_HOST")
		psqlPort := os.Getenv("POSTGRESQL_PORT")
		psqlUser := os.Getenv("POSTGRESQL_USER")
		psqlPassword := os.Getenv("POSTGRESQL_PASSWORD")
		psqldb := os.Getenv("POSTGRESQL_DB")
		repo, err := ps.NewPostgresRepository(psqlHost, psqlPort, psqlUser, psqlPassword, psqldb, timeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mysql":
		mysqlHost := os.Getenv("MYSQL_HOST")
		mysqlPort := os.Getenv("MYSQL_PORT")
		mysqlUser := os.Getenv("MYSQL_USER")
		mysqlPassword := os.Getenv("MYSQL_PASSWORD")
		mysqldb := os.Getenv("MYSQL_DB")
		repo, err := ms.NewMySQLRepository(mysqlHost, mysqlPort, mysqlUser, mysqlPassword, mysqldb, timeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}

	return nil
}
