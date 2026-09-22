# go-nem

<p align="center">
  <img src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_LightBlue.png" height="80" alt="Go">
  <br>
  <img src="assets/aemo.png" height="44" alt="AEMO">
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/shanehull/go-nem.svg)](https://pkg.go.dev/github.com/shanehull/go-nem)
[![Go Report Card](https://goreportcard.com/badge/github.com/shanehull/go-nem)](https://goreportcard.com/report/github.com/shanehull/go-nem)
[![Go](https://img.shields.io/github/go-mod/go-version/shanehull/go-nem)](https://go.dev/dl/)
[![CI](https://github.com/shanehull/go-nem/actions/workflows/test.yaml/badge.svg)](https://github.com/shanehull/go-nem/actions/workflows/test.yaml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Go client for AEMO National Electricity Market data published on [NEMWEB](https://nemweb.com.au). Standard library only, zero dependencies.

It discovers report files, downloads and caches them, and parses the two wire formats AEMO publishes: the nested record format used by MMS reports and the human-formatted text used by market notices.

## Install

```sh
go get github.com/shanehull/go-nem
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shanehull/go-nem"
)

func main() {
	client, err := nem.New(nem.WithCacheDir("/var/lib/nem"))
	if err != nil {
		log.Fatal(err)
	}

	notices, err := client.FetchNotices(context.Background(), nem.Since(time.Now().Add(-24*time.Hour)))
	if err != nil {
		log.Fatal(err)
	}
	for _, notice := range notices {
		fmt.Printf("%d  %-24s %s\n", notice.NoticeID, notice.TypeID, notice.IssueDate.Format("2006-01-02"))
	}
}
```

## Usage

### Listing files

`List` reads the report directory and returns the files matching the options, sorted by name, which orders them chronologically. Options filter by modification time, cap the count keeping the newest, and select the archive tier.

```go
refs, err := client.List(ctx, nem.ReportDispatchIS,
    nem.Since(time.Now().Add(-24*time.Hour)),
    nem.Limit(12),
)

// The archive tier, for a single day.
refs, err := client.List(ctx, nem.ReportMarketNotice,
    nem.FromArchive(),
    nem.On(time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC)),
)
```

### Market notices

Market notices are human-formatted text. `FetchNotices` lists and downloads the notices matching the options, and `DecodeMarketNotice` parses one file.

```go
notices, err := client.FetchNotices(ctx, nem.Limit(10))
for _, notice := range notices {
    fmt.Println(notice.NoticeID, notice.TypeID, notice.ReferencedIDs)
}
```

### MMS reports

MMS reports use the nested record format and are published as zip archives. `FetchTables` lists, downloads, and parses every matching file into `Table` values.

```go
tables, err := client.FetchTables(ctx, nem.ReportDispatchIS, nem.Limit(1))
for _, table := range tables {
    fmt.Println(table.Group, table.Name, table.Version, len(table.Rows))
}
```

### Decoding

Typed decoders read a `Table` slice into concrete structs.

```go
prices, err := nem.DecodeDispatchPrice(tables)
sums, err := nem.DecodeDispatchRegionSum(tables)
flows, err := nem.DecodeDispatchInterconnectorRes(tables)
fcas, err := nem.DecodeFCASPrice(tables)

units, err := nem.DecodeDispatchUnitScada(tables) // from ReportDispatchSCADA
```

### Caching

`WithCacheDir` enables an on-disk cache. Downloads are written atomically, and a cached file is never fetched twice.

```go
client, err := nem.New(
    nem.WithCacheDir("/var/lib/nem"),
    nem.WithMinFetchInterval(200*time.Millisecond),
)
```

### Time conventions

AEMO uses NEM time (AEST, UTC+10) with no daylight saving. MMS reports use year-first timestamps and market notices use day-first dates; the parsers detect the layout. Dispatch timestamps are interval-ending.

```go
t, err := nem.ParseNEMTime("2026/09/21 08:10:00") // MMS
t, err := nem.ParseNEMDate("26/08/2026")          // notice

day := nem.TradingDay(t)                 // "2026-09-21", the 04:05 to 04:00 trading day
start := nem.IntervalBeginning(t)        // 08:05, the interval that ends at 08:10
```

### Error handling

```go
refs, err := client.List(ctx, nem.ReportDispatchIS)
var apiErr nem.APIError
if errors.As(err, &apiErr) {
    fmt.Printf("NEMWEB error %d for %s\n", apiErr.StatusCode, apiErr.URL)
}
```

## API Coverage

| Area     | Functions                                                                                                                                                                                                                      |
| -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Client   | `New`, `WithBaseURL`, `WithHTTPClient`, `WithCacheDir`, `WithMinFetchInterval`, `WithUserAgent`                                                                                                                                |
| Listing  | `List`, `Since`, `Until`, `On`, `Limit`, `FromArchive`                                                                                                                                                                         |
| Fetching | `Download`, `FetchTables`, `FetchNotices`                                                                                                                                                                                      |
| Parsing  | `Parse`, `DecodeMarketNotice`                                                                                                                                                                                                  |
| Decoding | `DecodeDispatchPrice`, `DecodeDispatchRegionSum`, `DecodeDispatchInterconnectorRes`, `DecodeDispatchUnitScada`, `DecodeFCASPrice`                                                                                              |
| Cache    | `NewCache`, `Path`, `Get`, `Put`                                                                                                                                                                                               |
| Time     | `ParseNEMTime`, `ParseNEMDate`, `InNEM`, `TradingDay`, `IntervalBeginning`, `TradingIntervalBeginning`                                                                                                                         |
| Region   | `Region`, `Regions`, `ParseRegion`, `Region.Valid`, `Region.Name`                                                                                                                                                              |
| Reports  | `ReportMarketNotice`, `ReportDispatchIS`, `ReportDispatch`, `ReportDispatchSCADA`, `ReportP5`, `ReportNetwork`, `ReportMTPASADUIDAvailability`, `ReportPDPASADUIDAvailability`, `ReportSevenDayOutlookFull`, `ReportTradingIS` |

## Examples

See [examples/](examples/) for runnable programs:

- [marketnotice](examples/marketnotice/main.go) — list and decode recent market notices
- [dispatch](examples/dispatch/main.go) — fetch a dispatch interval and print regional prices
- [scada](examples/scada/main.go) — fetch a dispatch SCADA interval and print unit output

The package also carries godoc examples in [example_test.go](example_test.go) for listing, fetching, parsing, decoding, and the time helpers. Run an example with:

```sh
go run ./examples/marketnotice
```

## Testing

Unit tests run offline against golden fixtures in [testdata/](testdata/). Live NEMWEB tests are opt-in.

```sh
go test ./...
NEM_LIVE=1 go test ./...
```

## License

MIT. See [LICENSE](LICENSE).

The AEMO logo is a trademark of the Australian Energy Market Operator, used only to identify the data source this library reads.
