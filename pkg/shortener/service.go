package shortener

type RedirectService interface {
	Find(code string) (*Redirect, error)
	Store(redirectReq *Redirect) error
	Update(redirectReq *Redirect) error
}
