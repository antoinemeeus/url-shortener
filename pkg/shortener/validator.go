package shortener

// RedirectValidator interface for the port/adapter method to a repository
type RedirectValidator interface {
	ValidateNewCode(input *RedirectRequest) error
	ValidateURL(input *RedirectRequest) error
}
