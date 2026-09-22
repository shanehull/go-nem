package nem

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MarketNotice is a parsed market notice.
type MarketNotice struct {
	NoticeID          int
	TypeID            string
	TypeDescription   string
	ExternalReference string
	IssueDate         time.Time
	CreationDate      time.Time
	Body              string
	ReferencedIDs     []int
	Raw               string
}

var (
	referencedNoticeRE = regexp.MustCompile(`(?i)market notice\s+(\d+)`)
	noticeEndRE        = regexp.MustCompile(`^-{3,}\s*$`)
)

// DecodeMarketNotice parses a market notice text report.
func DecodeMarketNotice(data []byte) (MarketNotice, error) {
	text := string(data)
	n := MarketNotice{Raw: text}

	var bodyLines []string
	inBody := false
	for _, line := range strings.Split(text, "\n") {
		if inBody {
			trimmed := strings.TrimSpace(line)
			if noticeEndRE.MatchString(trimmed) || strings.EqualFold(trimmed, "END OF REPORT") {
				break
			}
			bodyLines = append(bodyLines, line)
			continue
		}

		if isFieldLine(line, "Reason") {
			inBody = true
			continue
		}

		key, value, ok := splitNoticeField(line)
		if !ok {
			continue
		}
		switch key {
		case "Notice ID":
			id, err := strconv.Atoi(value)
			if err != nil {
				return MarketNotice{}, fmt.Errorf("nem: notice id %q: %w", value, err)
			}
			n.NoticeID = id
		case "Notice Type ID":
			n.TypeID = value
		case "Notice Type Description":
			n.TypeDescription = value
		case "External Reference":
			n.ExternalReference = value
		case "Issue Date":
			t, err := ParseNEMDate(value)
			if err != nil {
				return MarketNotice{}, err
			}
			n.IssueDate = t
		case "Creation Date":
			t, err := ParseNEMTime(value)
			if err != nil {
				return MarketNotice{}, err
			}
			n.CreationDate = t
		}
	}

	n.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
	n.ReferencedIDs = referencedIDs(n.Body, n.NoticeID)
	return n, nil
}

func isFieldLine(line, key string) bool {
	trimmed := strings.TrimSpace(line)
	colon := strings.Index(trimmed, ":")
	if colon < 0 {
		return false
	}
	return strings.TrimSpace(trimmed[:colon]) == key
}

func splitNoticeField(line string) (string, string, bool) {
	colon := strings.Index(line, ":")
	if colon < 0 {
		return "", "", false
	}
	key := strings.TrimSpace(line[:colon])
	value := strings.TrimSpace(line[colon+1:])
	if key == "" || value == "" {
		return "", "", false
	}
	return key, value, true
}

func referencedIDs(body string, exclude int) []int {
	matches := referencedNoticeRE.FindAllStringSubmatch(body, -1)
	seen := make(map[int]bool)
	var ids []int
	for _, m := range matches {
		id, err := strconv.Atoi(m[1])
		if err != nil || id == exclude || seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	return ids
}
