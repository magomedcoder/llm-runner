package gen

import "embed"

//go:embed migrations/*.sql
var Postgres embed.FS
