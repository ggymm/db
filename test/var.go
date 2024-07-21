package test

import (
	"embed"
)

//go:embed create.sql
var CreateSQL string

//go:embed insert.sql
var InsertSQL string

//go:embed select.sql
var SelectSQL string

//go:embed index/*
var SelectIndexSQL embed.FS
