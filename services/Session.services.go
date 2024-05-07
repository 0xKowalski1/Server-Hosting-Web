package services

import (
	"0xKowalski1/server-hosting-web/config"

	"github.com/gorilla/sessions"
)

func NewSessionStore() *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(config.Envs.SessionStoreSecret))

	store.MaxAge(60 * 60 * 24 * 7) // Sessions valid for 7 days, might want to set this in config
	store.Options.Path = "/"
	store.Options.HttpOnly = config.Envs.SessionStoreHttpOnly
	store.Options.Secure = config.Envs.SessionStoreSecure

	return store
}
