package wikilinkui

import (
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/rs/zerolog"
)

const (
	logGroupName  = "WikiLinkLogs"
	logStreamName = "WikiLinkUI"
	regionName    = "us-west-2"
)

// LoggingHandler logs to disk and cloudwatch
type LoggingHandler struct {
	Logger  zerolog.Logger
	Session *session.Session
	Svc     *cloudwatchlogs.CloudWatchLogs
}

func MakeLoggingHandler(logLevel zerolog.Level, filewriter io.Writer) (*LoggingHandler, error) {
	var logHandler = &LoggingHandler{}
	var err error

	// Creating logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logwriter := io.MultiWriter(os.Stdout, filewriter, logHandler)
	logHandler.Logger = zerolog.New(logwriter).With().Str("service", "linkui").Timestamp().Logger().Level(logLevel)

	// Create session
	logHandler.Session, err = session.NewSession(&aws.Config{
		Region: aws.String(regionName),
	})
	if err != nil {
		return nil, err
	}

	// Create SVC
	logHandler.Svc = cloudwatchlogs.New(logHandler.Session)

	// Create log group
	createGroupInput := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(logGroupName),
	}
	logHandler.Svc.CreateLogGroup(createGroupInput)

	// Create stream
	createStreamInput := &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	}
	logHandler.Svc.CreateLogStream(createStreamInput)

	return logHandler, nil
}

func (log *LoggingHandler) Write(p []byte) (n int, err error) {
	logEvents := []*cloudwatchlogs.InputLogEvent{{
		Message:   aws.String(string(p)),
		Timestamp: aws.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
	}}

	putEventsInput := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		LogEvents:     logEvents,
	}

	log.Svc.PutLogEvents(putEventsInput)

	return len(p), nil
}
