package validator

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/antoinemeeus/url-shortener/pkg/shortener"
	errs "github.com/pkg/errors"
)

const (
	minRedirectCodeLength = 3
	maxRedirectCodeLength = 20
	alnumRegexString      = "^[a-zA-Z0-9]+$"
)

// RedirectValidator definition for the redirect input validator
type RedirectValidator struct{}

// ValidateNewCode validates new code format from user input
func (r *RedirectValidator) ValidateNewCode(input *shortener.RedirectRequest) error {
	if len(input.NewCode) <= minRedirectCodeLength {
		err := fmt.Errorf("new code value is too short. %d characters min", minRedirectCodeLength)
		return errs.Wrap(err, "service.Redirect.Validator")
	}

	if len(input.NewCode) > maxRedirectCodeLength {
		err := fmt.Errorf("new code value is too long. %d characters max", maxRedirectCodeLength)
		return errs.Wrap(err, "service.Redirect.Validator")
	}

	if !r.isAlphaNumeric(input.NewCode) {
		err := fmt.Errorf("new code value contains illegal characters. Only alphanumeric value accepted")
		return errs.Wrap(err, "service.Redirect.Validator")
	}

	return nil
}

// ValidateURL validates URL format from user input
func (r *RedirectValidator) ValidateURL(input *shortener.RedirectRequest) error {
	if !r.isValidURL(input.URL) {
		err := fmt.Errorf("invalid URL value")
		return errs.Wrap(err, "service.Redirect.Validator")
	}

	return nil
}

func (*RedirectValidator) isAlphaNumeric(value string) bool {
	return regexp.MustCompile(alnumRegexString).MatchString(value)
}

func (*RedirectValidator) isValidURL(value string) bool {
	var i int

	// checks needed as of Go 1.6 because of change https://github.com/golang/go/commit/617c93ce740c3c3cc28cdd1a0d712be183d0b328#diff-6c2d018290e298803c0c9419d8739885L195
	// emulate browser and strip the '#' suffix prior to validation. see issue-#237
	if i = strings.Index(value, "#"); i > -1 {
		value = value[:i]
	}

	if len(value) == 0 {
		return false
	}

	u, err := url.ParseRequestURI(value)

	if err != nil || u.Scheme == "" {
		return false
	}

	return true
}
