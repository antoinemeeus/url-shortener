package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
)

var (
	// ErrRedirectNotFound is returned when a redirect resource is not found
	ErrRedirectNotFound = errors.New("redirect Not Found")
	// ErrRedirectInvalid is returned when a redirect request is invalid. Used for when validation fails.
	ErrRedirectInvalid = errors.New("redirect Invalid")
	// ErrAlreadyExist is returned when there is a code is already in use and cannot be saved.
	ErrAlreadyExist = errors.New("new_code already exist")
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
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC()

	return r.redirectRepo.Store(redirect)
}

// Update validates a redirect update request, check if the new code doesn't already exist in database and saves it via the repository interface.
func (r *redirectService) Update(redirect *Redirect, newCode string) error {
	existingRedirect, err := r.redirectRepo.Find(redirect.Code)
	if err != nil {
		return errs.Wrap(ErrRedirectNotFound, err.Error())
	}

	redirectExistWithNewCode, _ := r.redirectRepo.Find(newCode)
	if redirectExistWithNewCode != nil {
		return errs.Wrap(ErrAlreadyExist, "service.Redirect.Update")
	}

	redirect.ID = existingRedirect.ID
	redirect.CreatedAt = existingRedirect.CreatedAt
	redirect.URL = existingRedirect.URL
	redirect.Code = newCode

	return r.redirectRepo.Store(redirect)
}
