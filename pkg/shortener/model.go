package shortener

// Redirect is the object model that is used for data transfer
type Redirect struct {
	Code      string `json:"code" bson:"code"`
	NewCode   string `json:"new_code" validate:"format=alnum & gte=0 & lte=15 | empty=true"`
	URL       string `json:"url" bson:"url" validate:"format=url & empty=false | empty=true"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}
