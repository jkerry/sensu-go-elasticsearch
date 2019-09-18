package eventprocessing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sensu/sensu-go/types"
)

var (
	stdin *os.File
)

func GetPipedEvent() (*types.Event, error) {
	if stdin == nil {
		stdin = os.Stdin
	}

	eventJSON, err := ioutil.ReadAll(stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin: %s", err.Error())
	}
	event := &types.Event{}
	err = json.Unmarshal(eventJSON, event)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stdin data: %s", err.Error())
	}

	if err = event.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate event: %s", err.Error())
	}

	return event, nil
}

type MetricValue struct {
	Timestamp string   `json:"timestamp"`
	Name      string   `json:"name"`
	Entity    string   `json:"entity"`
	Value     float64  `json:"value"`
	Namespace string   `json:"namespace"`
	Tags      []string `json:"tags"`
}

type EventValue struct {
	Timestamp string           `json:"timestamp"`
	Entity    *types.Entity    `json:"entity"`
	Check     *types.Check     `json:"check"`
	Metrics   *types.Metrics   `json:"namespace"`
	Metadata  types.ObjectMeta `json:"metadata"`
}

// {
// 	"name": "avg_cpu",
// 	"value": "56.0",
// 	"timestamp": "2019-03-30 12:30:00.45",
// 	"entity": "demo_test_agent",
// 	"namespace": "demo_jk185160",
// 	"tags": [
//    "company_jkte001",
//    "site_1001"
// 	]
// }

func parseTimestamp(timestamp int64) (string, error) {
	stringTimestamp := strconv.FormatInt(timestamp, 10)
	if len(stringTimestamp) > 10 {
		stringTimestamp = stringTimestamp[:10]
	}
	t, err := strconv.ParseInt(stringTimestamp, 10, 64)
	if err != nil {
		return "", err
	}
	return time.Unix(t, 0).Format(time.RFC3339), nil
}

func buildTag(key string, value string, prefix string) string {
	if len(prefix) > 0 {
		return fmt.Sprintf("%s_%s_%s", prefix, key, value)
	}
	return fmt.Sprintf("%s_%s", key, value)
}

func GetMetricFromPoint(point *types.MetricPoint, entityID string, namespaceID string, entityLabels map[string]string) (MetricValue, error) {
	var metric MetricValue

	metric.Entity = entityID
	metric.Namespace = namespaceID
	// Find metric name
	nameField := strings.Split(point.Name, ".")
	metric.Name = nameField[0]

	// Find metric timstamp
	unixTimestamp, err := parseTimestamp(point.Timestamp)
	if err != nil {
		return *new(MetricValue), fmt.Errorf("failed to validate event: %s", err.Error())
	}
	metric.Timestamp = unixTimestamp
	metric.Tags = make([]string, len(point.Tags)+len(entityLabels)+1)
	i := 0
	for _, tag := range point.Tags {
		metric.Tags[i] = buildTag(tag.Name, tag.Value, "")
		i++
	}
	for key, val := range entityLabels {
		metric.Tags[i] = buildTag(key, val, "entity")
		i++
	}
	metric.Tags[i] = fmt.Sprintf("sensu_entity_name_%s", entityID)
	metric.Value = point.Value
	return metric, nil
}

func ParseEventTimestamp(event *types.Event) (EventValue, error) {
	var eventValue EventValue

	eventValue.Entity = event.Entity
	eventValue.Check = event.Check
	eventValue.Metrics = event.Metrics
	eventValue.Metadata = event.ObjectMeta

	timestamp, err := parseTimestamp(event.Timestamp)
	if err != nil {
		return *new(EventValue), fmt.Errorf("failed to validate event: %s", err.Error())
	}

	eventValue.Timestamp = timestamp
	return eventValue, nil
}
