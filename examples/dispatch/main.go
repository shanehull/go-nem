// Command dispatch fetches a dispatch interval and prints regional prices.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shanehull/go-nem"
)

func main() {
	client, err := nem.New()
	if err != nil {
		log.Fatal(err)
	}

	tables, err := client.FetchTables(context.Background(), nem.ReportDispatchIS, nem.Limit(1))
	if err != nil {
		log.Fatal(err)
	}

	prices, err := nem.DecodeDispatchPrice(tables)
	if err != nil {
		log.Fatal(err)
	}
	for _, price := range prices {
		fmt.Printf("%s  %-5s  %9.2f\n", price.SettlementDate.In(nem.NEM).Format(time.RFC3339), price.RegionID, price.RRP)
	}
}
