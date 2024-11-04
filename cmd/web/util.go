package web

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed "assets"
var Files embed.FS

func ToString(value any) string {
	return fmt.Sprintf("%v", value)
}

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
