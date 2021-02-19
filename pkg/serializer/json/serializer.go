package json

import (
	"encoding/json"

	"github.com/antoinemeeus/url-shortener/pkg/shortener"
	"github.com/pkg/errors"
)

// Redirect definition for the json serializer
type Redirect struct{}

// Decode decodes json to shortener.Redirect object
func (r *Redirect) Decode(input []byte) (*shortener.RedirectRequest, error) {
	redirect := &shortener.RedirectRequest{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}

// Encode encodes shortener.Redirect object to json
func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	rr := shortener.RedirectResponse{
		Code: input.Code,
	}
	rawMsg, err := json.Marshal(rr)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return rawMsg, nil
}
