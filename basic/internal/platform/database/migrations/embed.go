// Package migrations holds the versioned SQL schema files. They are embedded
// into the binary at build time so a single executable can run them anywhere,
// with no loose .sql files to ship alongside it.
package migrations

import "embed"

//go:embed *.sql
var Files embed.FS
