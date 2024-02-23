package handlers

import (
	"log/slog"
	"sync"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// cookie name & flash messages key
// this should be a .env file.
const (
	sessionName          string = "fmessages"
	cookieStoreKeyLength int    = 32
)

var (
	cookieStore *sessions.CookieStore //nolint:gochecknoglobals //Singleton
	cookieOnce  sync.Once             //nolint:gochecknoglobals //Singleton
)

func getCookieStore() *sessions.CookieStore {
	cookieOnce.Do(func() {
		cookieStore = sessions.NewCookieStore(securecookie.GenerateRandomKey(cookieStoreKeyLength))
	})
	return cookieStore
}

// Set adds a new message to the cookie store.
func setFlashmessages(c echo.Context, kind, value string) {
	session, _ := getCookieStore().Get(c.Request(), sessionName)

	session.AddFlash(value, kind)

	if err := session.Save(c.Request(), c.Response()); err != nil {
		slog.ErrorContext(c.Request().Context(), "Error saving flash message", slog.Any("error", err))
	}
}

// Get receives flash messages from cookie store.
func getFlashmessages(c echo.Context, kind string) []string {
	session, _ := getCookieStore().Get(c.Request(), sessionName)

	fm := session.Flashes(kind)

	// if there are some messages…
	if len(fm) > 0 {
		if err := session.Save(c.Request(), c.Response()); err != nil {
			slog.ErrorContext(c.Request().Context(), "Error saving flash message", slog.Any("error", err))
		}
		// we start an empty strings slice that we
		// then return with messages
		var flashes []string
		for _, fl := range fm {
			// we add the messages to the slice
			flashes = append(flashes, fl.(string))
		}

		return flashes
	}

	return nil
}
