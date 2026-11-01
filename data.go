package passpasspass

import "embed"

// DefaultData holds the built-in convention cards and sub-conventions under data/.
//
//go:embed data
var DefaultData embed.FS
