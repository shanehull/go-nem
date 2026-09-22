package nem_test

import (
	"context"
	"fmt"
	"time"

	"github.com/shanehull/go-nem"
)

func ExampleClient_List() {
	client, _ := nem.New(nem.WithCacheDir("/var/lib/nem"))
	refs, _ := client.List(context.Background(), nem.ReportDispatchIS,
		nem.Since(time.Now().Add(-24*time.Hour)),
	)
	for _, ref := range refs {
		fmt.Printf("%s %d\n", ref.Name, ref.Size)
	}
}

func ExampleClient_FetchNotices() {
	client, _ := nem.New()
	notices, _ := client.FetchNotices(context.Background(),
		nem.Since(time.Now().Add(-24*time.Hour)),
	)
	for _, notice := range notices {
		fmt.Printf("%d %s\n", notice.NoticeID, notice.TypeID)
	}
}

func ExampleClient_FetchTables() {
	client, _ := nem.New()
	tables, _ := client.FetchTables(context.Background(), nem.ReportDispatchIS, nem.Limit(1))
	for _, table := range tables {
		fmt.Println(table.Group, table.Name, table.Version, len(table.Rows))
	}
}

func ExampleParse() {
	client, _ := nem.New()
	refs, _ := client.List(context.Background(), nem.ReportDispatchIS, nem.Limit(1))
	data, _ := client.Download(context.Background(), refs[0])
	tables, _ := nem.Parse(data)
	for _, table := range tables {
		fmt.Println(table.Group, table.Name, len(table.Rows))
	}
}

func ExampleDecodeMarketNotice() {
	client, _ := nem.New()
	refs, _ := client.List(context.Background(), nem.ReportMarketNotice, nem.Limit(1))
	data, _ := client.Download(context.Background(), refs[0])
	notice, _ := nem.DecodeMarketNotice(data)
	fmt.Println(notice.NoticeID)
}

func ExampleDecodeDispatchPrice() {
	client, _ := nem.New()
	tables, _ := client.FetchTables(context.Background(), nem.ReportDispatchIS, nem.Limit(1))
	prices, _ := nem.DecodeDispatchPrice(tables)
	for _, price := range prices {
		fmt.Printf("%s %s %.2f\n", price.SettlementDate.Format(time.RFC3339), price.RegionID, price.RRP)
	}
}

func ExampleDecodeDispatchRegionSum() {
	client, _ := nem.New()
	tables, _ := client.FetchTables(context.Background(), nem.ReportDispatchIS, nem.Limit(1))
	sums, _ := nem.DecodeDispatchRegionSum(tables)
	for _, sum := range sums {
		fmt.Printf("%s demand=%.1f\n", sum.RegionID, sum.TotalDemand)
	}
}

func ExampleParseNEMTime() {
	t, _ := nem.ParseNEMTime("2026/09/21 08:10:00")
	fmt.Println(t.In(nem.NEM).Format(time.RFC3339))
	// Output: 2026-09-21T08:10:00+10:00
}

func ExampleTradingDay() {
	t, _ := nem.ParseNEMTime("2026/09/21 03:30:00")
	fmt.Println(nem.TradingDay(t))
	// Output: 2026-09-20
}

func ExampleParseRegion() {
	region, _ := nem.ParseRegion("nsw1")
	fmt.Println(region, region.Name())
	// Output: NSW1 New South Wales
}
