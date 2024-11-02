package templates

import _ "embed"

//go:embed data/env.example
var EnvExample string

//go:embed data/main.go
var Main string

//go:embed data/go.mod.tmpl
var Gomod string

//go:embed data/go.sum.tmpl
var Gosum string

//go:embed data/query.sql
var Query string

//go:embed data/schema.sql
var Schema string

//go:embed data/sqlc.yaml
var Sqlc string

//go:embed data/handlers/greet.go
var Greet string

//go:embed data/routes/routes.go
var Routes string
