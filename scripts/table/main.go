package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/eqcp/db"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}

}

func run() error {

	if len(os.Args) < 2 {
		return fmt.Errorf("usage: %s <table>", os.Args[0])
	}

	arg := os.Args[1]
	err := db.Init(context.Background())
	if err != nil {
		return fmt.Errorf("db.Init: %w", err)
	}

	rows, err := db.Query(context.Background(), "DESCRIBE "+arg, nil)
	if err != nil {
		return fmt.Errorf("db.Query: %w", err)
	}

	arg = strings.TrimSuffix(arg, "s")

	w, err := os.Create(fmt.Sprintf("%s/table.go", arg))
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer w.Close()
	w.WriteString("package " + arg + "\n\n")
	w.WriteString("import (\n")
	w.WriteString("\t\"database/sql\"\n")
	w.WriteString("\t\"fmt\"\n")
	w.WriteString(")\n\n")

	w.WriteString("type Table struct {\n")

	for rows.Next() {
		var field, typeStr, null, key, extra string
		var defaultStr sql.NullString
		err = rows.Scan(&field, &typeStr, &null, &key, &defaultStr, &extra)
		if err != nil {
			return fmt.Errorf("rows.Scan: %w", err)
		}

		if strings.Contains(typeStr, "int") {
			typeStr = "int"
			if null == "YES" {
				typeStr = "sql.NullInt64"
			}
		} else if strings.Contains(typeStr, "varchar") {
			typeStr = "string"
			if null == "YES" {
				typeStr = "sql.NullString"
			}
		} else if strings.Contains(typeStr, "text") {
			typeStr = "string"
			if null == "YES" {
				typeStr = "sql.NullString"
			}
		} else if strings.Contains(typeStr, "tinyint") {

			typeStr = "bool"
			if null == "YES" {
				typeStr = "sql.NullBool"
			}
		} else if strings.Contains(typeStr, "datetime") {
			typeStr = "time.Time"
			if null == "YES" {
				typeStr = "sql.NullTime"
			}
		} else if strings.Contains(typeStr, "float") {
			typeStr = "float64"
			if null == "YES" {
				typeStr = "sql.NullFloat64"
			}
		} else {
			return fmt.Errorf("unknown type: %s", typeStr)
		}

		fieldName := strings.ToLower(field)
		if fieldName == "id" {
			fieldName = "ID"
		}
		fieldName = strings.ToUpper(fieldName[:1]) + fieldName[1:]

		w.WriteString("\t" + fieldName + " " + typeStr + " `db:\"" + field + "\"`\n")
		fmt.Printf("Field: %s, Type: %s, Null: %s, Key: %s, Default: %+v, Extra: %s\n", field, typeStr, null, key, defaultStr, extra)
	}

	w.WriteString("}\n\n")
	fmt.Println("Generated", arg+"/table.go")

	return nil
}
