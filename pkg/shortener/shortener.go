package shortener

import (
	"errors"
	"fmt"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	// ErrRedirectNotFound is returned when a redirect resource is not found
	ErrRedirectNotFound = errors.New("redirect Not Found")
	// ErrRedirectInvalid is returned when a redirect request is invalid. Used for when validation fails.
	ErrRedirectInvalid = errors.New("redirect Invalid")
	// ErrAlreadyExist is returned when there is a code is already in use and cannot be saved.
	ErrAlreadyExist = errors.New("code already exist")
	// ErrNewCodeEmpty is returned when the newCode provided is empty.
	ErrNewCodeEmpty = errors.New("new_code is empty")
)

type redirectService struct {
	redirectRepo RedirectRepository
}

// NewRedirectService returns a new instance of the redirectService.
func NewRedirectService(redirectRepo RedirectRepository) RedirectService {
	return &redirectService{
		redirectRepo,
	}
}

// Find returns a redirect resource via the repository interface
func (r *redirectService) Find(code string) (*Redirect, error) {
	redirect, err := r.redirectRepo.Find(code)
	if err != nil {
		return nil, errs.Wrap(ErrRedirectNotFound, err.Error())
	}
	return redirect, nil
}

// Store validates a redirect creation request and saves it via the repository interface.
func (r *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, fmt.Sprintf("service.Redirect.Store Validation Error: %s", err.Error()))
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC()
	return r.redirectRepo.Store(redirect)
}

// Update validates a redirect update request, check if the new code doesn't already exist in database and saves it via the repository interface.
func (r *redirectService) Update(redirect *Redirect) error {
	if redirect.NewCode == "" {
		return errs.Wrap(ErrNewCodeEmpty, "service.Redirect.Store")
	}

	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, fmt.Sprintf("service.Redirect.Store Validation Error: %s", err.Error()))
	}

	if redirectExistWithNewCode, _ := r.redirectRepo.Find(redirect.NewCode); redirectExistWithNewCode != nil {
		return errs.Wrap(ErrAlreadyExist, "service.Redirect.Update")
	}

	oldRedirect, err := r.redirectRepo.Find(redirect.Code)
	if err != nil {
		return errs.Wrap(ErrRedirectNotFound, err.Error())
	}

	redirect.Code = redirect.NewCode
	redirect.URL = oldRedirect.URL
	redirect.CreatedAt = time.Now().UTC()
	_ = r.redirectRepo.Delete(oldRedirect) // Explicitly ignore error, but needs to be logged in the future to find memory leaks in repo

	return r.redirectRepo.Store(redirect)
}
