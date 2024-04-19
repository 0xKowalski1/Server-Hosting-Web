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
	store := sessions.NewCookieStore([]byte("my-secret")) // take me from config

	store.MaxAge(60 * 60 * 24 * 7) // Sessions valid for 7 days
	store.Options.Path = "/"
	store.Options.HttpOnly = false // Take me from config
	store.Options.Secure = false   // me too

	return store
}
