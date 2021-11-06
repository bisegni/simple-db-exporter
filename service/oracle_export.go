package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/sijms/go-ora"
)

// ExportFromFolderToFolder ...
func ExportFromFolderToFolder(
	oracle_uri string,
	fromFolderPath string,
	toFolderPath string,
) error {
	exists, err := FSExists(fromFolderPath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("From folder not exists")
	}

	exists, err = FSExists(toFolderPath)
	if err != nil {
		return err
	}
	if !exists {
		if err := FSCreateDirectory(toFolderPath); err != nil {
			return err
		}
	}
	// open connection
	conn, err := sql.Open("oracle", fmt.Sprintf("oracle://%s", oracle_uri))
	if err != nil {
		return err
	}

	if err := FSWorkOnItemInPath(
		fromFolderPath,
		func(entry os.FileInfo) error {
			//skip dir
			if entry.IsDir() {
				return nil
			}
			sourceFile := filepath.Join(
				fromFolderPath,
				entry.Name(),
			)

			destFile := filepath.Join(
				toFolderPath,
				fmt.Sprintf(
					"%s-insert.sql",
					strings.TrimSuffix(
						entry.Name(),
						filepath.Ext(entry.Name()),
					),
				),
			)
			fmt.Printf("Processing file %s into %s\n", sourceFile, destFile)
			// process file
			return exportSingleFile(
				conn,
				sourceFile,
				destFile,
			)
		},
	); err != nil {
		return err
	}
	return nil
}

func exportSingleFile(
	conn *sql.DB,
	sqlFilePath string,
	sqlInsertOutputFile string,
) error {
	exists, err := FSExists(sqlInsertOutputFile)
	if err != nil {
		return err
	}
	if exists {
		if err := FSRemove(sqlInsertOutputFile); err != nil {
			return nil
		}
	}
	f, err := CreateFile(sqlInsertOutputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// check for err
	rows, err := conn.Query(string(content))
	if err != nil {
		return err
	}

	col, err := rows.Columns()
	if err != nil {
		return err
	}
	readCols := make([]interface{}, len(col))
	writeCols := make([]interface{}, len(col))
	for i, _ := range writeCols {
		readCols[i] = &writeCols[i]
	}
	destTableName := strings.TrimSuffix(
		filepath.Base(sqlFilePath),
		filepath.Ext(filepath.Base(sqlFilePath)),
	)
	f.WriteString(fmt.Sprintf("TRUNCATE TABLE %s;\n", destTableName))
	f.WriteString("BEGIN;\n")
	for rows.Next() {
		// define vars
		err := rows.Scan(readCols...)
		// check for error
		if err != nil {
			return err
		}

		var b bytes.Buffer
		for i := 0; i < len(col); i++ {
			writeDataInBuffer(writeCols[i], &b)

			if i+1 < len(col) {
				b.WriteRune(',')
			}
		}

		f.WriteString(fmt.Sprintf("INSERT INTO %s values (%s);\n", destTableName, b.String()))
	}
	f.WriteString("COMMIT;\n")
	// check for error
	defer rows.Close()
	return nil
}

// ExportOracleQuery ...
func ExportOracleQuery(
	oracle_uri string,
	sqlFilePath string,
	destTableName string,
	destinationFilePath string,
) error {
	destFilePathName := filepath.Join(destinationFilePath, fmt.Sprintf("%s-insert.sql", destTableName))
	if err := FSCreateDirectory(destinationFilePath); err != nil {
		return nil
	}

	exists, err := FSExists(sqlFilePath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("sql file not exists")
	}

	exists, err = FSExists(destFilePathName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("destination file already exists")
	}

	f, err := CreateFile(destFilePathName)
	if err != nil {
		return err
	}
	defer f.Close()

	conn, err := sql.Open("oracle", fmt.Sprintf("oracle://%s", oracle_uri))
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// check for err
	rows, err := conn.Query(string(content))
	if err != nil {
		return err
	}

	col, err := rows.Columns()
	if err != nil {
		return err
	}
	readCols := make([]interface{}, len(col))
	writeCols := make([]interface{}, len(col))
	for i, _ := range writeCols {
		readCols[i] = &writeCols[i]
	}

	f.WriteString(fmt.Sprintf("TRUNCATE TABLE %s;\n", destTableName))
	f.WriteString("BEGIN;\n")
	for rows.Next() {
		// define vars
		err := rows.Scan(readCols...)
		// check for error
		if err != nil {
			return err
		}

		var b bytes.Buffer
		for i := 0; i < len(col); i++ {
			writeDataInBuffer(writeCols[i], &b)

			if i+1 < len(col) {
				b.WriteRune(',')
			}
		}

		f.WriteString(fmt.Sprintf("INSERT INTO %s values (%s);\n", destTableName, b.String()))
	}
	f.WriteString("COMMIT;\n")
	// check for error
	defer rows.Close()

	// check for err
	defer func() {
		err := conn.Close()
		// check for err
		if err != nil {
			fmt.Println(err)
		}
	}()
	return nil
}

func writeDataInBuffer(column interface{}, buf *bytes.Buffer) {
	switch v := column.(type) {
	case nil:
		buf.WriteString("null")
	case string:
		buf.WriteString(
			fmt.Sprintf("'%v'", escapeString(v)),
		)
	default:
		// And here I'm feeling dumb. ;)
		buf.WriteString(fmt.Sprintf("%v", v))
	}
}

func escapeString(value string) string {
	var sb strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch c {
		case '\\', 0, '\n', '\r', '\'', '"':
			sb.WriteByte('\\')
			sb.WriteByte(c)
		case '\032':
			sb.WriteByte('\\')
			sb.WriteByte('Z')
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}
