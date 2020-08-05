package postgresql

import (
	"fmt"

	"github.com/antoinemeeus/url-shortener/pkg/shortener"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type postgresRepository struct {
	database *gorm.DB
}

// NewPostgresRepository returns a new instance of the PostgresQL repository.
func NewPostgresRepository(host string, port string, user string, password string, dbName string, timeout int) (shortener.RedirectRepository, error) {
	repo := &postgresRepository{}
	args := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s connect_timeout=%d sslmode=disable", host, port, user, dbName, password, timeout)
	db, err := gorm.Open("postgres", args)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	db.AutoMigrate(&shortener.Redirect{})
	repo.database = db
	return repo, nil
}

// Find finds the corresponding URL for the code provided and construct the shortener.Redirect object from saved information.
func (r *postgresRepository) Find(code string) (*shortener.Redirect, error) {
	sr := &shortener.Redirect{}
	err := r.database.Where(&shortener.Redirect{Code: code}).First(sr).Error
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	return sr, nil
}

// Store stores or update a new code and URL to PostgresQL via a ORM from the shortener.Redirect object.
func (r *postgresRepository) Store(redirect *shortener.Redirect) error {
	var err error
	if r.database.NewRecord(*redirect) {
		err = r.database.Create(redirect).Error
	}
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}

// Delete deletes a shortener.Redirect entry by record
func (r *postgresRepository) Delete(redirect *shortener.Redirect) error {
	err := r.database.Unscoped().Delete(redirect).Error
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Delete")
	}
	return nil
}

// Close allow to close database connection gracefully
func (r *postgresRepository) Close() error {
	return r.database.Close()
}
