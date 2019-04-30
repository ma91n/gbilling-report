package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

const QUERY = `
SELECT
  project.id AS pj_name,
  service.description AS sdict,
  SUM(cost) AS cost
FROM
  ` + "`dev-tky-dx-233223.billing.gcp_billing_export_v1_01FAB2_66C27C_15BAB5`" +`
WHERE
  FORMAT_DATE('%Y-%m-%d',DATE_SUB(CURRENT_DATE(),INTERVAL 2 DAY)) = FORMAT_DATE('%Y-%m-%d', DATE(export_time))
GROUP BY
  project.id,
  service.description
ORDER BY
  project.id,
  service.description
LIMIT
  100
`

func main() {

	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	if len(projectID) == 0 {
		log.Fatal("env GOOGLE_PROJECT_ID is not found")
	}
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	it, err := client.Query(QUERY).Read(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var values []bigquery.Value
		if err := it.Next(&values); err == iterator.Done {
			break
		}

		if err != nil {
			log.Printf("Failed to Iterate Query:%v\n", err)
		}

		log.Print(values)
	}

}
