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

// maxZipDepth bounds how many nested zip layers Parse unwraps. AEMO archives
// some reports double-zipped, and a corrupt file should not loop forever.
const maxZipDepth = 4

// errNoTables reports a stream carrying no nested-record tables, so a zip with
// non-report entries can skip them.
var errNoTables = errors.New("nem: no tables found")

// Parse decodes a nested record report into its tables. It accepts raw nested
// CSV bytes or a zip archive, unwrapping nested zips and every entry of a
// multi-entry archive. AEMO publishes daily dispatch reports as one outer zip
// containing one inner zip per interval.
func Parse(data []byte) ([]Table, error) {
	return parse(data, 0)
}

func parse(data []byte, depth int) ([]Table, error) {
	if !bytes.HasPrefix(data, zipMagic) {
		return parseNested(data)
	}
	if depth >= maxZipDepth {
		return nil, errors.New("nem: too many nested zip layers")
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("nem: open zip: %w", err)
	}

	var all []Table
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		body, err := readZipEntry(f)
		if err != nil {
			return nil, err
		}
		tables, err := parse(body, depth+1)
		if errors.Is(err, errNoTables) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("nem: parse zip entry %q: %w", f.Name, err)
		}
		all = append(all, tables...)
	}
	if len(all) == 0 {
		return nil, errNoTables
	}
	return all, nil
}

func readZipEntry(f *zip.File) ([]byte, error) {
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
			return nil, fmt.Errorf("nem: unknown record type %q", clip(rec[0], 32))
		}
	}

	if len(tables) == 0 {
		return nil, errNoTables
	}
	return tables, nil
}

func tableKey(group, name string, version int) string {
	return group + "\x00" + name + "\x00" + strconv.Itoa(version)
}

// clip shortens a value for an error message so a binary blob does not produce
// a multi-kilobyte error line.
func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
