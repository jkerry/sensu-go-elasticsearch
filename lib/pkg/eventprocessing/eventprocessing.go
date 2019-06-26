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

	if !event.HasMetrics() {
		return nil, fmt.Errorf("event does not contain metrics")
	}

	return event, nil
}

type Tag struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}
type MetricValue struct {
	Timestamp string  `json:"timestamp"`
	Name      string  `json:"name"`
	Entity    string  `json:"entity"`
	Value     float64 `json:"value"`
	Namespace string  `json:"namespace"`
	Tags      []Tag   `json:"tags"`
}

// {
// 	"name": "avg_cpu",
// 	"value": "56.0",
// 	"timestamp": "2019-03-30 12:30:00.45",
// 	"entity": "demo_test_agent",
// 	"namespace": "demo_jk185160",
// 	"tags": [
// 			{
// 					"tag": "company",
// 					"value": "JKTE001"
// 			},
// 			{
// 					"tag": "site",
// 					"value": "1001"
// 			}
// 	]
// }

func parsePointTimestamp(point *types.MetricPoint) (string, error) {
	stringTimestamp := strconv.FormatInt(point.Timestamp, 10)
	if len(stringTimestamp) > 10 {
		stringTimestamp = stringTimestamp[:10]
	}
	t, err := strconv.ParseInt(stringTimestamp, 10, 64)
	if err != nil {
		return "", err
	}
	return time.Unix(t, 0).Format(time.RFC3339), nil
}

func buildTag(key string, value string, prefix string) Tag {
	var tag Tag
	var tagName string
	if len(prefix) > 0 {
		tagName = fmt.Sprintf("%s_%s", prefix, key)
	} else {
		tagName = key
	}
	tag.Tag = tagName
	tag.Value = value
	return tag
}

func GetMetricFromPoint(point *types.MetricPoint, entityID string, namespaceID string, entityLabels map[string]string) (MetricValue, error) {
	var metric MetricValue

	metric.Entity = entityID
	metric.Namespace = namespaceID
	// Find metric name
	nameField := strings.Split(point.Name, ".")
	metric.Name = nameField[0]

	// Find metric timstamp
	unixTimestamp, err := parsePointTimestamp(point)
	if err != nil {
		return *new(MetricValue), fmt.Errorf("failed to validate event: %s", err.Error())
	}
	metric.Timestamp = unixTimestamp
	metric.Tags = make([]Tag, len(point.Tags)+len(entityLabels)+1)
	i := 0
	for _, tag := range point.Tags {
		metric.Tags[i] = buildTag(tag.Name, tag.Value, "")
		i++
	}
	for key, val := range entityLabels {
		metric.Tags[i] = buildTag(key, val, "entity_label")
		i++
	}
	var entityNameTag Tag
	entityNameTag.Tag = "sensu_entity_name"
	entityNameTag.Value = entityID
	metric.Tags[i] = entityNameTag
	metric.Value = point.Value
	return metric, nil
}