package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/antoinemeeus/url-shortener/pkg/shortener"
	errs "github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type sqlRepository struct {
	database *gorm.DB
}

// NewMySQLRepository returns a new instance of the PostgresQL repository.
func NewMySQLRepository(host string, port string, user string, password string, dbName string, timeout int) (shortener.RedirectRepository, error) {
	repo := &sqlRepository{}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewMySQLRepository")
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	db.WithContext(timeoutContext)
	_ = db.AutoMigrate(&shortener.Redirect{})
	repo.database = db

	return repo, nil
}

// Find finds the corresponding URL for the code provided and construct the shortener.Redirect object from saved information.
func (r *sqlRepository) Find(code string) (*shortener.Redirect, error) {
	sr := &shortener.Redirect{}
	err := r.database.Where(&shortener.Redirect{Code: code}).First(sr).Error
	if err != nil {
		return nil, errs.Wrap(err, "repository.Redirect.Find")
	}

	return sr, nil
}

// Store stores or update a new code and URL to MySQL via a ORM from the shortener.Redirect object.
func (r *sqlRepository) Store(redirect *shortener.Redirect) error {
	var err error

	err = r.database.First(&shortener.Redirect{},redirect.ID).Error
	if err != nil {
		err = r.database.Create(redirect).Error
		if err != nil {
			return errs.Wrap(err, "repository.Redirect.Store")
		}
		return nil
	}
	err = r.database.Model(redirect).Update("code",redirect.Code).Error
	if err != nil {
		return errs.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}

// Delete deletes a shortener.Redirect entry by record
func (r *sqlRepository) Delete(redirect *shortener.Redirect) error {
	err := r.database.Unscoped().Delete(redirect).Error
	if err != nil {
		return errs.Wrap(err, "repository.Redirect.Delete")
	}

	return nil
}

// Close allow to close database connection gracefully
func (r *sqlRepository) Close() error {
	db, _ := r.database.DB()

	return db.Close()
}
