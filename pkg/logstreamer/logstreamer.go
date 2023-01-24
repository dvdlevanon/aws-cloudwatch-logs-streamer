package logstreamer

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"k8s.io/utils/pointer"
)

type LogStreamer struct {
	interval     int
	logs         *cloudwatchlogs.CloudWatchLogs
	logGroupName string
}

func New(logGroupName string, interval int) (*LogStreamer, error) {
	s, err := session.NewSession()

	if err != nil {
		return nil, err
	}

	return &LogStreamer{
		logs:         cloudwatchlogs.New(s),
		logGroupName: logGroupName,
		interval:     interval,
	}, nil
}

func (r *LogStreamer) GetLogStreams() ([]string, error) {
	streams, err := r.logs.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &r.logGroupName,
	})

	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, stream := range streams.LogStreams {
		result = append(result, *stream.LogStreamName)
	}

	return result, nil
}

func (r *LogStreamer) Read(logStreamName string, events chan<- string) {
	var next *string
	startMillis := time.Now().UnixMilli()

	fmt.Fprintf(os.Stderr, "Streaming logs from group name: %s, stream: %s\n", r.logGroupName, logStreamName)
	for {
		time.Sleep(time.Duration(r.interval) * time.Millisecond)
		result, err := r.logs.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  &r.logGroupName,
			LogStreamName: &logStreamName,
			StartTime:     &startMillis,
			NextToken:     next,
			StartFromHead: pointer.Bool(true),
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting logs from stream: %s - %v\n", logStreamName, err)
			continue
		}

		for _, event := range result.Events {
			events <- *event.Message
		}

		next = result.NextForwardToken
	}
}
