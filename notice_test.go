package nem_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/shanehull/go-nem"
)

func TestDecodeMarketNotice(t *testing.T) {
	notice, err := nem.DecodeMarketNotice(readFixture(t, "market_notice.txt"))
	if err != nil {
		t.Fatalf("DecodeMarketNotice: %v", err)
	}
	if notice.NoticeID != 144924 {
		t.Errorf("NoticeID = %d, want 144924", notice.NoticeID)
	}
	if notice.TypeID != "MINIMUM SYSTEM LOAD" {
		t.Errorf("TypeID = %q", notice.TypeID)
	}
	if notice.TypeDescription != "MSL1/MSL2/MSL3" {
		t.Errorf("TypeDescription = %q", notice.TypeDescription)
	}
	if !notice.IssueDate.Equal(mustTime(t, "2026-08-26T00:00:00+10:00")) {
		t.Errorf("IssueDate = %s", notice.IssueDate)
	}
	if !notice.CreationDate.Equal(mustTime(t, "2026-08-26T11:18:33+10:00")) {
		t.Errorf("CreationDate = %s", notice.CreationDate)
	}
	if notice.ExternalReference == "" {
		t.Error("ExternalReference is empty")
	}
	if !strings.Contains(notice.Body, "AEMO ELECTRICITY MARKET NOTICE") {
		t.Errorf("Body missing notice header: %q", notice.Body)
	}
	if strings.Contains(notice.Body, "END OF REPORT") {
		t.Errorf("Body should stop before the report footer: %q", notice.Body)
	}
}

func TestDecodeMarketNoticeReferencedIDs(t *testing.T) {
	notice, err := nem.DecodeMarketNotice(readFixture(t, "market_notice.txt"))
	if err != nil {
		t.Fatalf("DecodeMarketNotice: %v", err)
	}
	want := []int{144917, 144808}
	if !reflect.DeepEqual(notice.ReferencedIDs, want) {
		t.Errorf("ReferencedIDs = %v, want %v", notice.ReferencedIDs, want)
	}
}
