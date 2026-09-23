# AGENTS.md — go-nem

AEMO NEM data client for Go. Stdlib only. No dependencies.

## Commands

```bash
go build ./...          # always builds clean (no deps)
go vet ./...            # always clean
go test ./...           # unit tests, offline, no network
NEM_LIVE=1 go test ./...  # adds live NEMWEB tests
go test ./... -run TestParseDispatchIS   # run a single test
```

Offline tests parse golden fixtures in `testdata/`. Live tests are gated on `NEM_LIVE` and skip cleanly when unset.

## Architecture

```
nem.go        Client, New(), WithBaseURL/WithHTTPClient/WithCacheDir/WithMinFetchInterval/WithUserAgent
errors.go     type APIError = internal.APIError  (type alias, use errors.As)
types.go      Report, Tier, Kind, FileRef, Table
reports.go    Report definitions (directory and file prefix)
options.go    ListOption types: Since, Until, On, Limit, FromArchive
list.go       List, directory-listing parse, filtering
fetch.go      Download, FetchTables, FetchNotices, request pacing
parse.go      Nested record parser (I, D, C rows), zip detection
notice.go     MarketNotice and DecodeMarketNotice (text report parser)
mmsdm.go      Typed decoders over Table: price, region sum, interconnector, SCADA, FCAS
time.go       NEM time, date parse, trading day, interval helpers
region.go     NEM market regions: Region, Regions, ParseRegion
cache.go      Cache: path mapping, atomic writes
internal/http.go  Get — all HTTP I/O
```

Every public method takes `ctx context.Context` as first argument.

## Gotchas

### Two wire formats

MMS reports use the nested record format: `I` rows define a table header, `D` rows are data, `C` rows are comments. Market notices are human-formatted text, not the nested format, and are parsed by `DecodeMarketNotice`.

### Two date formats

MMS `SETTLEMENTDATE` is year-first (`2026/09/21 08:10:00`). Notice dates are day-first (`26/08/2026`). `ParseNEMTime` and `ParseNEMDate` detect the layout from the first field. Notice fields pad with runs of spaces; the parsers collapse whitespace.

### Nested record layout

An `I` row is `I,<group>,<table>,<version>,<column...>`. A `D` row is `D,<group>,<table>,<version>,<value...>`. Columns and values both start at index 4. A `D` row is attached to the `I` row with the same group, table, and version. The verified reports carry no record count in the header. Values are comma-separated and quoted when they contain a space, so the parser uses `encoding/csv` with `LazyQuotes`.

### Report files

Dispatch and MMS reports are zip archives containing one CSV. Market notices are plain text. `Parse` detects zip by magic bytes (`PK\x03\x04`) and otherwise treats the input as raw CSV.

### Directory listings

NEMWEB serves IIS directory listings, not Apache. The whole listing is one line, `HREF` is uppercase, and the modification time and size precede the anchor (`Monday, September 21, 2026 08:05 AM   20629 <A HREF=...>`). `List` matches anchors by the report's file prefix, so the `Ancillary_Services_Payments` report (which does not share a `PUBLIC_` prefix) needs its prefix confirmed before it is added.

### Time conventions

AEMO uses NEM time (AEST, UTC+10, no daylight saving). Dispatch timestamps are interval-ending. A trading day runs 04:05 to 04:00 the following day. All conversions live in `time.go`.

## Scope

This library covers NEM market data: NEMWEB report discovery and parsing, market notices, dispatch and price tables, and the market region set.

## Conventions

- Tests in the `nem_test` external package. Test functions named `Test<Name>`.
- When adding a decoder, also add a test, an example if it is user-facing, and a row in the README API coverage table.
- Live tests use `os.Getenv("NEM_LIVE")` and `t.Skip` when unset.
- Semantic commits for release-please. `feat:` bumps minor (pre-1.0), `fix:` bumps patch.
