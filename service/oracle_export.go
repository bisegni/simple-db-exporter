package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/sijms/go-ora/v2"
)

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

		f.WriteString(fmt.Sprintf("INSERT INTO %s values (%s)\n", destTableName, b.String()))
	}
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
	case string:
		buf.WriteString(fmt.Sprintf("'%v'", v))
	default:
		// And here I'm feeling dumb. ;)
		buf.WriteString(fmt.Sprintf("%v", v))
	}
}
