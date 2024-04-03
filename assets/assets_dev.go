//go:build dev
// +build dev

package assets

import (
	"net/http"
)

func Assets() http.FileSystem {
	return http.Dir("assets")
}
