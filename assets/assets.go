package assets

import (
	"embed"
	"net/http"
)

//go:embed fonts/* img/* js/* css/*
var assetsFS embed.FS

func Assets() http.FileSystem {
	return http.FS(assetsFS)
}
