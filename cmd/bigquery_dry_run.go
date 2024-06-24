package cmd

import (
	"context"
	"io"
	"log"
	"time"
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
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

    for i:=0; i < 2; i++ {
        ctx, cancel := context.WithTimeout(ctx, 6 * time.Second)
        defer cancel()

        job, err := q.Run(ctx)

        if err == nil {
            status := job.LastStatus()
            if err := status.Err(); err != nil {
                return 0., err
            }
            bytes_processed := float32(status.Statistics.TotalBytesProcessed)
            return bytes_processed, err
        }

        if i==1 {
            return 0.,err
        }
    }
    return 0., err
}
