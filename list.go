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

var listingLineRE = regexp.MustCompile(`href="([^"]+)"[^>]*>[^<]*</a>\s*(\d{2}-[A-Za-z]{3}-\d{4})?\s*(\d{2}:\d{2})?\s*(\d+)?`)

// List returns the files published under a report directory, sorted by name,
// which orders them chronologically. It lists the current tier by default.
func (c *Client) List(ctx context.Context, report Report, opts ...ListOption) ([]FileRef, error) {
	p := &listParams{tier: TierCurrent}
	for _, opt := range opts {
		opt(p)
	}

	dirURL := c.baseURL + "/Reports/" + p.tier.String() + "/" + report.Dir + "/"
	body, err := internal.Get(ctx, c.httpClient, dirURL, c.userAgent)
	if err != nil {
		return nil, err
	}

	return filterRefs(parseListing(string(body), report, p.tier, c.baseURL), p), nil
}

func parseListing(body string, report Report, tier Tier, baseURL string) []FileRef {
	var refs []FileRef
	for _, line := range strings.Split(body, "\n") {
		m := listingLineRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		href := m[1]
		name := path.Base(href)
		if strings.HasSuffix(href, "/") || !strings.HasPrefix(name, report.Prefix) {
			continue
		}
		ref := FileRef{
			Report: report,
			Name:   name,
			URL:    baseURL + "/Reports/" + tier.String() + "/" + report.Dir + "/" + url.PathEscape(name),
			Tier:   tier,
		}
		if m[2] != "" && m[3] != "" {
			if t, err := time.ParseInLocation("02-Jan-2006 15:04", m[2]+" "+m[3], NEM); err == nil {
				ref.Modified = t
			}
		}
		if m[4] != "" {
			if size, err := strconv.ParseInt(m[4], 10, 64); err == nil {
				ref.Size = size
			}
		}
		refs = append(refs, ref)
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].Name < refs[j].Name })
	return refs
}

func filterRefs(refs []FileRef, p *listParams) []FileRef {
	out := make([]FileRef, 0, len(refs))
	for _, r := range refs {
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
