package billing

import (
	"log"
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"strconv"
)

type bqRepository struct {
	client *bigquery.Client
	tableName string
}

func New(projectID, tableName string) (*bqRepository, error) {
	client, err := bigquery.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	return &bqRepository{client: client, tableName: tableName}, nil
}

func (r *bqRepository) FetchCost(termDay int) (TotalCost, error) {

	query := `
SELECT
  sub.project_name,
  sub.export_time,
  sub.cost,
  sub.currency
FROM (
  SELECT
    project.id AS project_name,
    FORMAT_TIMESTAMP('%Y-%m-%d', export_time) AS export_time,
    ROUND(SUM(cost)) AS cost,
    MAX(currency)    as currency
  FROM
    ` + "`" + r.tableName  +"`" + `
  WHERE
    export_time >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(),INTERVAL ` + strconv.Itoa(termDay) + ` DAY)
  AND project.id is not null
  GROUP BY
    project.id,
    FORMAT_TIMESTAMP('%Y-%m-%d', export_time) ) sub
ORDER BY
  sub.project_name,
  sub.export_time
`

	it, err := r.client.Query(query).Read(context.Background())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var timeseries []Cost
	for {
		var c Cost
		err := it.Next(&c)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate query:%v\n", err)
			break
		}
		timeseries = append(timeseries, c)
	}
	return timeseries, nil
}
