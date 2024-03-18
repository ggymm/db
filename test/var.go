package test

import _ "embed"

//go:embed create.sql
var CreateSQL string

//go:embed insert.sql
var InsertSQL string

//go:embed select.sql
var SelectSQL string
