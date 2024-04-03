//go:build prod
// +build prod

package assets

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed fonts/* img/* js/* css/*
var assetsFS embed.FS

func Assets() http.FileSystem {
	// print the contents of the assets
	entries, err := assetsFS.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil
	}

	for _, name := range entries {
		fmt.Println(name.Name())

		if name.IsDir() {
			subEntries, err := assetsFS.ReadDir(name.Name())
			if err != nil {
				fmt.Println("Error reading subdirectory:", err)
				return nil
			}

			for _, subName := range subEntries {
				fmt.Println("  ", subName.Name())
			}
		}
	}

	return http.FS(assetsFS)
}
