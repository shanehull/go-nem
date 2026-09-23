package nem

import (
	"context"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shanehull/go-nem/internal"
)

// NEMWEB serves IIS directory listings: the whole listing is one line, HREF is
// uppercase, and the modification time and size precede the anchor.
var (
	anchorRE = regexp.MustCompile(`(?i)<a\s+href="([^"]+)"[^>]*>([^<]*)</a>`)
	iisMeta  = regexp.MustCompile(`([A-Za-z]+,\s+[A-Za-z]+\s+\d{1,2},\s+\d{4}\s+\d{1,2}:\d{2}\s+[AP]M)\s+(\d+)\s*$`)
)

// List returns the files published under a report directory, sorted by name,
// which orders them chronologically. It lists the current tier by default.
func (c *Client) List(ctx context.Context, report Report, opts ...ListOption) ([]FileRef, error) {
	p := &listParams{tier: TierCurrent}
	for _, opt := range opts {
		opt(p)
	}

	// No trailing slash: NEMWEB's host answers ".../<Report>/" with
	// "Website is unavailable", while ".../<Report>" returns the listing.
	dirURL := c.baseURL + "/Reports/" + p.tier.String() + "/" + report.Dir
	body, err := internal.Get(ctx, c.httpClient, dirURL, c.userAgent)
	if err != nil {
		return nil, err
	}

	return filterRefs(parseListing(string(body), report, p.tier, c.baseURL), p), nil
}

func parseListing(body string, report Report, tier Tier, baseURL string) []FileRef {
	matches := anchorRE.FindAllStringSubmatchIndex(body, -1)
	refs := make([]FileRef, 0, len(matches))
	for i, loc := range matches {
		href := body[loc[2]:loc[3]]
		name := path.Base(href)
		if strings.HasSuffix(href, "/") || !strings.HasPrefix(name, report.Prefix) {
			continue
		}

		beforeStart := 0
		if i > 0 {
			beforeStart = matches[i-1][1]
		}
		modified, size := parseListingMeta(body[beforeStart:loc[0]])

		refs = append(refs, FileRef{
			Report:   report,
			Name:     name,
			URL:      baseURL + "/Reports/" + tier.String() + "/" + report.Dir + "/" + url.PathEscape(name),
			Size:     size,
			Modified: modified,
			Tier:     tier,
		})
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].Name < refs[j].Name })
	return refs
}

// parseListingMeta extracts the modification time and size that a listing puts
// either before the anchor (IIS) or after it (Nginx and Apache). A miss returns
// the zero value, which callers treat as unknown rather than filtering on it.
func filterRefs(refs []FileRef, p *listParams) []FileRef {
	out := make([]FileRef, 0, len(refs))
	for _, r := range refs {
		if p.nameStamp != "" && !strings.Contains(r.Name, p.nameStamp) {
			continue
		}
		if !r.Modified.IsZero() {
			if !p.since.IsZero() && r.Modified.Before(p.since) {
				continue
			}
			if !p.until.IsZero() && !r.Modified.Before(p.until) {
				continue
			}
		}
		out = append(out, r)
	}
	if p.limit > 0 && len(out) > p.limit {
		out = out[len(out)-p.limit:]
	}
	return out
}

// parseListingMeta extracts the IIS modification time and size that precede an
// anchor. A miss returns the zero value, which callers treat as unknown rather
// than filtering on it.
func parseListingMeta(before string) (time.Time, int64) {
	m := iisMeta.FindStringSubmatch(strings.TrimSpace(before))
	if m == nil {
		return time.Time{}, 0
	}
	t, err := time.ParseInLocation("Monday, January 2, 2006 3:04 PM", m[1], NEM)
	if err != nil {
		return time.Time{}, 0
	}
	size, err := strconv.ParseInt(m[2], 10, 64)
	if err != nil {
		return time.Time{}, 0
	}
	return t, size
}
