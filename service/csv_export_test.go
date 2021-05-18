package service_test

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"testing"

	"github.com/bisegni/simple-db-exporter/service"
	"github.com/stretchr/testify/assert"
)

func TestSelectAllCol(t *testing.T) {
	var filePath = "test_data.csv"
	os.RemoveAll(filePath)
	defer os.RemoveAll(filePath)

	err := service.RunCSVExport(
		"postgres://postgres:postgres@localhost:5432/test-db",
		"test_data",
		[]string{},
		-1,
		filePath,
	)
	assert.NoError(t, err)

	row, err := lineCounter(filePath)
	assert.NoError(t, err)
	assert.True(t, row == 300)

	col, err := getColumnNumber(filePath)
	assert.NoError(t, err)
	assert.True(t, col == 2)
}

func TestSelectFilteredColumn(t *testing.T) {
	var filePath = "test_data.csv"
	os.RemoveAll(filePath)
	defer os.RemoveAll(filePath)

	err := service.RunCSVExport(
		"postgres://postgres:postgres@localhost:5432/test-db",
		"test_data",
		[]string{"name"},
		-1,
		filePath,
	)
	assert.NoError(t, err)

	row, err := lineCounter(filePath)
	assert.NoError(t, err)
	assert.True(t, row == 300)

	col, err := getColumnNumber(filePath)
	assert.NoError(t, err)
	assert.True(t, col == 1)
}

func lineCounter(filePath string) (int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := f.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func getColumnNumber(filePath string) (int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	record, err := csvReader.Read()
	if err != nil {
		return 0, err
	}

	return len(record), nil
}
