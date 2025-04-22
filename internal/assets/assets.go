package assets

//go:generate pnpm build

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed *.js *.js.map
var assets embed.FS

func FS() fs.FS {
	return assets
}

func Handler() http.Handler {
	return http.FileServerFS(assets)
}
