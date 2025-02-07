// Пакет содержит миграции в двух форматах - в виде sql-скриптов и
// в виде go-файлов. Оба формата миграций при компиляции встраиваются в
// исполняемый файл.
// go install github.com/pressly/goose/v3/cmd/goose@latest
// goose create add_some_table sql
// goose create add_some_table go
// goose --help

package migrations

import "embed"

//go:embed *.sql
var EmbedMigrations embed.FS
