package json

import (
	"testing"
	"time"

	"github.com/antoinemeeus/url-shortener/pkg/shortener"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRedirect_Decode(t *testing.T) {
	re := &Redirect{}
	jsn := `{"code":"firstCode","new_code":"newCode123","url":"/newCode"}`
	actual, err := re.Decode([]byte(jsn))
	expected := &shortener.RedirectRequest{
		Code:    "firstCode",
		URL:     "/newCode",
		NewCode: "newCode123",
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
		Model: shortener.Model{
			ID:        1,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Code: "firstCode",
		URL:  "/newCode",
	}

	actual, err := re.Encode(payload)
	expected := []byte(`{"code":"firstCode"}`)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
