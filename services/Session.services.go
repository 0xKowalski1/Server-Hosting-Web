package services

import (
	"github.com/gorilla/sessions"
)

// Config
// CookiesKey
// MaxAge
// HttpOnly
// Secure

func NewSessionStore() *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte("my-secret"))

	store.MaxAge(60 * 60 * 24 * 2)
	store.Options.Path = "/"
	store.Options.HttpOnly = false
	store.Options.Secure = false

	return store
}
