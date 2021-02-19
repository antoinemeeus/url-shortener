package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/antoinemeeus/url-shortener/pkg/shortener"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewRedisRepository returns a new instance of the Redis repository.
func NewRedisRepository(redisURL string) (shortener.RedirectRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	repo.client = client

	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

// Find finds the corresponding URL for the code provided and construct the shortener.Redirect object from saved information.
func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = time.Unix(createdAt, 0)
	return redirect, nil
}

// Store stores a new code and URL to Redis from the shortener.Redirect object.
func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt.Unix(),
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}

// Delete delete a short.Redirect entry in Redis
func (r *redisRepository) Delete(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	err := r.client.Del(key).Err()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Delete")
	}

	return nil
}

// Close Allow to close connection gracefully
func (r *redisRepository) Close() {
	_ = r.client.Close()
}
