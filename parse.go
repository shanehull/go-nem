package nem

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var zipMagic = []byte{'P', 'K', 0x03, 0x04}

// Parse decodes a nested record report into its tables. It accepts either a
// zip archive containing one CSV or raw nested CSV bytes.
func Parse(data []byte) ([]Table, error) {
	if bytes.HasPrefix(data, zipMagic) {
		raw, err := firstZipEntry(data)
		if err != nil {
			return nil, err
		}
		data = raw
	}
	return parseNested(data)
}

func firstZipEntry(data []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("nem: open zip: %w", err)
	}
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("nem: open zip entry %q: %w", f.Name, err)
		}
		body, readErr := io.ReadAll(rc)
		closeErr := rc.Close()
		if readErr != nil {
			return nil, fmt.Errorf("nem: read zip entry %q: %w", f.Name, readErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("nem: close zip entry %q: %w", f.Name, closeErr)
		}
		return body, nil
	}
	return nil, errors.New("nem: zip archive is empty")
}

func parseNested(data []byte) ([]Table, error) {
	r := csv.NewReader(bytes.NewReader(data))
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	var tables []Table
	index := make(map[string]int)

	for {
		rec, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("nem: read record: %w", err)
		}
		if len(rec) == 0 {
			continue
		}

		switch strings.TrimSpace(rec[0]) {
		case "C":
			continue
		case "I":
			if len(rec) < 5 {
				return nil, fmt.Errorf("nem: header row has %d fields, want at least 5", len(rec))
			}
			version, err := strconv.Atoi(rec[3])
			if err != nil {
				return nil, fmt.Errorf("nem: header %s/%s version %q: %w", rec[1], rec[2], rec[3], err)
			}
			tables = append(tables, Table{
				Group:   rec[1],
				Name:    rec[2],
				Version: version,
				Columns: rec[4:],
			})
			index[tableKey(rec[1], rec[2], version)] = len(tables) - 1
		case "D":
			if len(rec) < 5 {
				return nil, fmt.Errorf("nem: data row has %d fields, want at least 5", len(rec))
			}
			version, err := strconv.Atoi(rec[3])
			if err != nil {
				return nil, fmt.Errorf("nem: data %s/%s version %q: %w", rec[1], rec[2], rec[3], err)
			}
			i, ok := index[tableKey(rec[1], rec[2], version)]
			if !ok {
				return nil, fmt.Errorf("nem: data row for unknown table %s/%s v%d", rec[1], rec[2], version)
			}
			tables[i].Rows = append(tables[i].Rows, rec[4:])
		default:
			return nil, fmt.Errorf("nem: unknown record type %q", rec[0])
		}
	}

	if len(tables) == 0 {
		return nil, errors.New("nem: no tables found")
	}
	return tables, nil
}

func tableKey(group, name string, version int) string {
	return group + "\x00" + name + "\x00" + strconv.Itoa(version)
}
