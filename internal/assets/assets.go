package assets

//go:generate pnpm build

import (
	"embed"
	"io/fs"
)

//go:embed *.js *.js.map
var assets embed.FS

func FS() fs.FS {
	return assets
}
