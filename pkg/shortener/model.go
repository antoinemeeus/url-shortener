package shortener

import (
	"time"
)

// gorm.Model definition
type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Redirect is the object model that is used for data transfer
type Redirect struct {
	Model
	Code    string `json:"code" gorm:"type:varchar(100);unique_index"`
	NewCode string `json:"new_code" gorm:"type:varchar(100)" validate:"format=alnum & gte=0 & lte=15 | empty=true"`
	URL     string `json:"url" validate:"format=url & empty=false | empty=true"`
}
