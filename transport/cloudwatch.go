package transport

import (
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/google/uuid"
	"github.com/sheacloud/goflow-addons/utils"
)

type CloudwatchState struct {
	LogGroupName      string
	LogStream         *LogStream
	CloudwatchLogsSvc *cloudwatchlogs.CloudWatchLogs
}

func (c *CloudwatchState) Initialize() {
	logStreamName := uuid.New().String()

	c.LogStream = &LogStream{
		LogStreamName:       logStreamName,
		LogGroupName:        c.LogGroupName,
		bufferedEvents:      []*cloudwatchlogs.InputLogEvent{},
		bufferedEventsBytes: 0,
		cloudwatchState:     c,
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
	LogStreamName       string
	LogGroupName        string
	LastSequenceToken   *string
	bufferedEvents      []*cloudwatchlogs.InputLogEvent
	bufferedEventsBytes int
	bufferedEventLock   sync.Mutex
	lastUploadTime      time.Time
	uploadLock          sync.Mutex
	cloudwatchState     *CloudwatchState
}

func (l *LogStream) IngestEvents(msgs []*utils.ExtendedFlowMessage) {
	for _, msg := range msgs {
		data, _ := utils.HumanReadableJSONMarshal(msg)
		l.addEventToBuffer(data, int64(msg.TimeReceived)*1000)
	}
}

func (l *LogStream) addEventToBuffer(message []byte, timestamp int64) {
	event := &cloudwatchlogs.InputLogEvent{
		Message:   aws.String(string(message)),
		Timestamp: aws.Int64(timestamp),
	}
	l.bufferedEventLock.Lock()
	l.bufferedEvents = append(l.bufferedEvents, event)
	l.bufferedEventsBytes += len(message) + 26

	// fmt.Printf("Added %v byte event to buffer, size is now %v bytes\n", len(message)+26, l.bufferedEventsBytes)

	var uploadEvents bool

	// TODO figure out better buffer size cutoff to correspond to AWS limit of 1,048,576 bytes
	if l.bufferedEventsBytes >= 500000 {
		uploadEvents = true
		fmt.Println("Uploading events as 100KB buffer has been reached")
	} else if time.Now().Sub(l.lastUploadTime).Seconds() >= 10 {
		uploadEvents = true
		fmt.Println("Uploading events as time limit has been reached")
	} else if len(l.bufferedEvents) > 9000 {
		uploadEvents = true
		fmt.Println("Uploading events as 9k event limit has been reached")
	}
	l.bufferedEventLock.Unlock()

	if uploadEvents {
		l.UploadBufferedEvents()
	}
}

func (l *LogStream) UploadBufferedEvents() {
	// acquire the buffered event lock, copy the data out of the list, and reset it
	l.bufferedEventLock.Lock()
	eventList := make([]*cloudwatchlogs.InputLogEvent, len(l.bufferedEvents))
	copy(eventList, l.bufferedEvents)
	l.bufferedEvents = nil
	l.bufferedEventsBytes = 0
	l.bufferedEventLock.Unlock()

	l.uploadLock.Lock()
	resp, err := l.cloudwatchState.CloudwatchLogsSvc.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
		LogEvents:     eventList,
		LogGroupName:  aws.String(l.LogGroupName),
		LogStreamName: aws.String(l.LogStreamName),
		SequenceToken: l.LastSequenceToken,
	})
	if err != nil {
		fmt.Println("error publishing log events")
		fmt.Println(err.Error())
		l.uploadLock.Unlock()
		return
	}

	l.LastSequenceToken = resp.NextSequenceToken
	fmt.Println(resp.RejectedLogEventsInfo)
	l.lastUploadTime = time.Now()
	l.uploadLock.Unlock()
}

func (c *CloudwatchState) Publish(msgs []*utils.ExtendedFlowMessage) {
	c.LogStream.IngestEvents(msgs)
}
