package db

import "embed"

//go:embed schema
var schemaFiles embed.FS

func GetSchema() ([]byte, error) {
	return schemaFiles.ReadFile("schema/schema.sql")
}
