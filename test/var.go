package test

import _ "embed"

//go:embed create.sql
var CreateSQL string

//go:embed sql_list.txt
var SQLList string