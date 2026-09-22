// Command marketnotice lists and decodes recent market notices.
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

	notices, err := client.FetchNotices(context.Background(), nem.Since(time.Now().Add(-24*time.Hour)))
	if err != nil {
		log.Fatal(err)
	}
	for _, notice := range notices {
		fmt.Printf("%d  %-24s %s\n", notice.NoticeID, notice.TypeID, notice.IssueDate.Format("2006-01-02"))
	}
}
