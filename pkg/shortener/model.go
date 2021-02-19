package shortener

import (
	"time"

	"gorm.io/gorm"
)

// Model gorm.Model definition
type Model struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// Redirect is the object model that is used for data transfer
type Redirect struct {
	Model
	Code    string `json:"code" gorm:"type:varchar(100);unique_index"`
	URL     string `json:"url"`
}

// RedirectResponse is the object used to return the redirect response
type RedirectResponse struct {
	Code string `json:"code"`
}

// RedirectRequest is the object used to return the redirect response
type RedirectRequest struct {
	URL     string `json:"url"`
	Code    string `json:"code"`
	NewCode string `json:"new_code"`
}
