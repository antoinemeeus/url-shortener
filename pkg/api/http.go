package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	js "github.com/antoinemeeus/url-shortener/pkg/serializer/json"
	"github.com/antoinemeeus/url-shortener/pkg/shortener"
	"github.com/antoinemeeus/url-shortener/pkg/validator"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// RedirectHandler interface that defines the http controller api methods
type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService shortener.RedirectService
}

// NewHandler creates a new instance of the handler containing the redirectService.
func NewHandler(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{redirectService: redirectService}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) serializer(contentType string) shortener.RedirectSerializer {
	if contentType == "application/json" {
		return &js.Redirect{}
	}

	return &js.Redirect{}
}

func (h *handler) validator() shortener.RedirectValidator {
	return &validator.RedirectValidator{}
}

// Get controller action for GET requests. Returns the corresponding url for a existing code.
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}

// Post controller action for POST requests. Saves a new url and return its corresponding code.
func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	redirectInput, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.validator().ValidateURL(redirectInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	redirect := shortener.Redirect{
		URL:  redirectInput.URL,
	}
	err = h.redirectService.Store(&redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseBody, err := h.serializer(contentType).Encode(&redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(w, contentType, responseBody, http.StatusCreated)
}

// Update controller action for PUT requests. Updates the corresponding code to a new personalized code.
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	redirectInput, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = h.validator().ValidateNewCode(redirectInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	redirect := shortener.Redirect{
		Code: redirectInput.Code,
	}
	err = h.redirectService.Update(&redirect, redirectInput.NewCode)
	if err != nil {
		switch errors.Cause(err) {
		case shortener.ErrRedirectInvalid:
			http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusBadRequest), err.Error()), http.StatusBadRequest)

		case shortener.ErrNewCodeEmpty:
			http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusBadRequest), err.Error()), http.StatusBadRequest)

		case shortener.ErrRedirectNotFound:
			http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusNotFound), err.Error()), http.StatusNotFound)

		case shortener.ErrAlreadyExist:
			http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusForbidden), err.Error()), http.StatusForbidden)

		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	responseBody, err := h.serializer(contentType).Encode(&redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(w, contentType, responseBody, http.StatusCreated)
}
