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
func RunCSVExport(uri string, table string, extractColumn []string, file string) error {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
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

	for rows.Next() {
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

	return nil
}
