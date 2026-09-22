// Command scada fetches a dispatch SCADA interval and prints unit output.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shanehull/go-nem"
)

const limit = 20

func main() {
	client, err := nem.New()
	if err != nil {
		log.Fatal(err)
	}

	tables, err := client.FetchTables(context.Background(), nem.ReportDispatchSCADA, nem.Limit(1))
	if err != nil {
		log.Fatal(err)
	}

	units, err := nem.DecodeDispatchUnitScada(tables)
	if err != nil {
		log.Fatal(err)
	}

	for i, unit := range units {
		if i >= limit {
			fmt.Printf("... %d more units\n", len(units)-limit)
			break
		}
		fmt.Printf("%s %-12s %9.2f MW\n", unit.SettlementDate.In(nem.NEM).Format("15:04"), unit.DUID, unit.SCADAValue)
	}
}
