package json

import (
	"testing"

	"github.com/antoinemeeus/url-shortener/pkg/shortener"
	"github.com/stretchr/testify/assert"
)

func TestRedirect_Decode(t *testing.T) {
	re := &Redirect{}
	jsn := `{"code":"firstCode","new_code":"newCode","url":"/newCode"}`
	actual, err := re.Decode([]byte(jsn))
	expected := &shortener.Redirect{
		Code:    "firstCode",
		NewCode: "newCode",
		URL:     "/newCode",
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestRedirect_Decode_With_Bad_Format_Will_Error(t *testing.T) {
	re := &Redirect{}
	jsn := `{`
	actual, err := re.Decode([]byte(jsn))

	assert.Error(t, err)
	assert.EqualError(t, err, "serializer.Redirect.Decode: unexpected end of JSON input")
	assert.Nil(t, actual)
}

func TestRedirect_Encode(t *testing.T) {
	re := &Redirect{}
	payload := &shortener.Redirect{
		Code:    "firstCode",
		NewCode: "newCode",
		URL:     "/newCode",
	}
	actual, err := re.Encode(payload)
	expected := `{"code":"firstCode","new_code":"newCode","url":"/newCode"}`

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
