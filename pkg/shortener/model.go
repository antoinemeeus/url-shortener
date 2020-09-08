package shortener

import (
	"time"
)

// Model gorm.Model definition
type Model struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Redirect is the object model that is used for data transfer
type Redirect struct {
	Model
	Code    string `json:"code" gorm:"type:varchar(100);unique_index"`
	NewCode string `json:"new_code" gorm:"type:varchar(100)"validate:"eq=0 | empty=false & format=alnum & gte=3 & lte=20"`
	URL     string `json:"url" validate:"format=url & empty=false | empty=true"`
}
