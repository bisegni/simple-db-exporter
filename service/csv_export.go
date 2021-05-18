package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/jackc/pgx/v4"
)

// RunCSVExport ...
func RunCSVExport(uri string, table string, extractColumn []string, rowLimit int, file string) error {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	var totalRecord int = 0
	fmt.Printf("Connecting to %s ", uri)
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())
	fmt.Println("DONE!")

	rows, err := conn.Query(context.Background(), fmt.Sprintf("select * from %s", table))
	if err != nil {
		return fmt.Errorf("QueryRow failed: %v", err)
	}
	defer rows.Close()

	if err := os.MkdirAll(path.Dir(file), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	csvW := csv.NewWriter(f)
	defer csvW.Flush()
	//now scan the data

	fieldDescriptions := rows.FieldDescriptions()
	var columns []string
	for _, col := range fieldDescriptions {
		columns = append(columns, string(col.Name))
	}

	sort.Strings(extractColumn)
	fmt.Print("exporting ")
	for rows.Next() {
		if rowLimit != -1 &&
			totalRecord >= rowLimit {
			break
		}
		if totalRecord%10 == 0 {
			fmt.Print(".")
		}
		totalRecord++
		var entry []string
		data, err := rows.Values()
		if err != nil {
			return err
		}

		for i, col := range columns {
			if len(extractColumn) > 0 {
				l := sort.SearchStrings(extractColumn, col)
				if l == len(extractColumn) || extractColumn[l] == col {
					continue
				}
			}
			entry = append(entry, fmt.Sprintf("%v", data[i]))
		}
		if err := csvW.Write(entry); err != nil {
			return err
		}
	}
	fmt.Println(" done!")
	return nil
}
