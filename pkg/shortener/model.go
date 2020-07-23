package shortener

type Redirect struct {
	Slug      string `json:"slug" bson:"slug" validate:"empty=false & lte=10 & format=alnum_unicode"`
	URL       string `json:"url" bson:"url" validate:"empty=false & format=url"`
	CreatedAt int64  `json:"created_at" bson:"created_at" msgpack:"created_at"`
}
