package cmd

import (
	"cloud.google.com/go/bigquery"
	"context"
	"google.golang.org/api/option"
	"io"
	"log"
)

func createBigQueryClient(ctx context.Context, projectId string, keyfile string) (*bigquery.Client, error) {
	if keyfile == "" {
		return bigquery.NewClient(ctx, projectId)
	} else {
		return bigquery.NewClient(ctx, projectId, option.WithCredentialsFile(keyfile))
	}
}

func queryDryRun(w io.Writer, query *string, projectId string, keyfile string, location string) (float32, error) {
	ctx := context.Background()

	var client *bigquery.Client
	var err error

	client, err = createBigQueryClient(ctx, projectId, keyfile)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	q := client.Query(*query)
	q.DryRun = true
	q.DisableQueryCache = false
	q.Location = location

	job, err := q.Run(ctx)
	if err != nil {
		return 0., err
	}

	status := job.LastStatus()
	if err := status.Err(); err != nil {
		return 0., err
	}
	bytes_processed := float32(status.Statistics.TotalBytesProcessed)
	return bytes_processed, err

}
