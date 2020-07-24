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
	ErrRedirectNotFound = errors.New("redirect Not Found")
	ErrRedirectInvalid  = errors.New("redirect Invalid")
	ErrAlreadyExist     = errors.New("code already exist")
)

type redirectService struct {
	redirectRepo RedirectRepository
}

func NewRedirectService(redirectRepo RedirectRepository) RedirectService {
	return &redirectService{
		redirectRepo,
	}
}

func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.redirectRepo.Find(code)
}

func (r *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, fmt.Sprintf("service.Redirect.Store Validation Error: %s", err.Error()))
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepo.Store(redirect)
}

func (r *redirectService) Update(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, fmt.Sprintf("service.Redirect.Store Validation Error: %s", err.Error()))
	}

	if redirectExistWithNewCode, _ := r.redirectRepo.Find(redirect.NewCode); redirectExistWithNewCode != nil {
		return errs.Wrap(ErrAlreadyExist, "service.Redirect.Update")
	}

	oldRedirect, err := r.redirectRepo.Find(redirect.Code)
	if err != nil {
		return errs.Wrap(ErrRedirectNotFound, fmt.Sprintf("service.Redirect.Update Error: %s", err.Error()))
	}

	redirect.Code = redirect.NewCode
	redirect.URL = oldRedirect.URL
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepo.Store(redirect)
}
