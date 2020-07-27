package shortener

// RedirectSerializer interface for the port/adapter for serialization and deserialization
type RedirectSerializer interface {
	Decode(input []byte) (*Redirect, error)
	Encode(input *Redirect) ([]byte, error)
}
