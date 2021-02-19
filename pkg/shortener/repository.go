package shortener

// RedirectRepository interface for the port/adapter method to a repository
type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
	Delete(redirect *Redirect) error
	Close()
}
