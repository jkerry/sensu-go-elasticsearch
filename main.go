package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/jkerry/sensu-go-elasticsearch/lib/pkg/eventprocessing"
	"github.com/spf13/cobra"
)

var (
	index         string
	dated_postfix bool
)

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sensu-go-elasticsearch",
		Short: "The Sensu Go handler for metric and event logging in elasticsearch\nRequired:  Set the ELASTICSEARCH_URL env var with an appropriate connection url (https://user:pass@hostname:port)",
		RunE:  run,
	}

	cmd.Flags().BoolVarP(&dated_postfix,
		"dated_index",
		"d",
		false,
		"Should the index have the current date postfixed? ie: metric_data-2019-06-27")

	cmd.Flags().StringVarP(&index,
		"index",
		"i",
		"",
		"metric_data")
	_ = cmd.MarkFlagRequired("index")
	return cmd
}

func generateIndex() string {
	if dated_postfix {
		dt := time.Now()
		return fmt.Sprintf("%s-%s", index, dt.Format("2006.01.02"))
	}
	return index
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	event, err := eventprocessing.GetPipedEvent()
	if err != nil {
		fmt.Errorf("Could not process or validate event data from stdin: %v", err)
		return err
	}

	for _, point := range event.Metrics.Points {
		metric, err := eventprocessing.GetMetricFromPoint(point, event.Entity.Name, event.Entity.Namespace, event.Entity.Labels)
		if err != nil {
			fmt.Errorf("error processing sensu event MetricPoints into MetricValue: %v", err)
			return err
		}
		msg, err := json.Marshal(metric)
		if err != nil {
			fmt.Errorf("error serializing metric data to pub/sub json payload: %v", err)
			return err
		}
		err = sendElasticSearchMetric(string(msg), index)
		if err != nil {
			fmt.Printf("error sending metric data to elasticsearch: %v", err)
			return err
		}
	}
	return nil
}

func sendElasticSearchMetric(metricBody string, index string) error {
	es, _ := elasticsearch.NewDefaultClient()
	req := esapi.IndexRequest{
		Index:   generateIndex(),
		Body:    strings.NewReader(metricBody),
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return fmt.Errorf("Error getting response: %s", err)
	}
	if res.IsError() {
		return fmt.Errorf("[%s] Error indexing document ID=%d", res.Status(), 0)
	}
	return nil
}
