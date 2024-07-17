package test

import (
	"embed"
	_ "embed"
)

//go:embed create.sql
var CreateSQL string

//go:embed insert.sql
var InsertSQL string

//go:embed update.sql
var UpdateSQL string

//go:embed delete.sql
var DeleteSQL string

//go:embed select.sql
var SelectSQL string

//go:embed select_where.sql
var SelectWhereSQL string

//go:embed index/*
var SelectIndexSQL embed.FS
