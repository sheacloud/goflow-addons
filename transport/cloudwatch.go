package transport

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	flowmessage "github.com/cloudflare/goflow/v3/pb"
	"github.com/google/uuid"
)

type CloudwatchState struct {
	LogGroupName      string
	LogStream         *LogStream
	CloudwatchLogsSvc *cloudwatchlogs.CloudWatchLogs
}

func (c *CloudwatchState) Initialize() {
	logStreamName := uuid.New().String()

	c.LogStream = &LogStream{
		LogStreamName: logStreamName,
		LogGroupName:  c.LogGroupName,
	}

	_, err := c.CloudwatchLogsSvc.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(c.LogGroupName),
		LogStreamName: aws.String(logStreamName),
	})
	if err != nil {
		fmt.Println("error creating log stream")
		return
	}
}

type LogStream struct {
	LogStreamName     string
	LogGroupName      string
	LastSequenceToken *string
}

func (l *LogStream) Publish(msgs []*flowmessage.FlowMessage, cloudwatchLogsSvc *cloudwatchlogs.CloudWatchLogs) {
	logEvents := []*cloudwatchlogs.InputLogEvent{}
	for _, msg := range msgs {
		data, _ := HumanReadableJSONMarshal(msg)
		logEvents = append(logEvents, &cloudwatchlogs.InputLogEvent{
			Message:   aws.String(string(data)),
			Timestamp: aws.Int64(int64(msg.TimeReceived) * 1000),
		})
		fmt.Printf("event size: %v bytes\n", len(data)+26)
	}

	resp, err := cloudwatchLogsSvc.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
		LogEvents:     logEvents,
		LogGroupName:  aws.String(l.LogGroupName),
		LogStreamName: aws.String(l.LogStreamName),
		SequenceToken: l.LastSequenceToken,
	})
	if err != nil {
		fmt.Println("error publishing log events")
		fmt.Println(err.Error())
		return
	}

	l.LastSequenceToken = resp.NextSequenceToken
}

func (c *CloudwatchState) Publish(msgs []*flowmessage.FlowMessage) {
	c.LogStream.Publish(msgs, c.CloudwatchLogsSvc)
}
