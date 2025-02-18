package log

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"rakamin-final-task/helpers/appcontext"
)

type LoggerLib struct {
	log *logrus.Logger
}

type LogInterface interface {
	Debug(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Fatal(context.Context, string, ...interface{})
}

func Init() LogInterface {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fileInfo(8)
		},
	})
	log.SetReportCaller(true)

	return &LoggerLib{
		log: log,
	}
}

func (l *LoggerLib) Debug(ctx context.Context, message string, field ...interface{}) {
	l.log.WithContext(ctx).WithFields(getFields(ctx, field...)).Debug(message)
}

func (l *LoggerLib) Info(ctx context.Context, message string, field ...interface{}) {
	l.log.WithContext(ctx).WithFields(getFields(ctx, field...)).Info(message)
}

func (l *LoggerLib) Warn(ctx context.Context, message string, field ...interface{}) {
	l.log.WithContext(ctx).WithFields(getFields(ctx, field...)).Warn(message)
}

func (l *LoggerLib) Error(ctx context.Context, message string, field ...interface{}) {
	l.log.WithContext(ctx).WithFields(getFields(ctx, field...)).Error(message)
}

func (l *LoggerLib) Fatal(ctx context.Context, message string, field ...interface{}) {
	l.log.WithContext(ctx).WithFields(getFields(ctx, field...)).Fatal(message)
}

func getFields(ctx context.Context, fields ...interface{}) logrus.Fields {
	logFields := logrus.Fields{
		"request_id":      appcontext.GetRequestId(ctx),
		"service_version": appcontext.GetServiceVersion(ctx),
		"user_agent":      appcontext.GetUserAgent(ctx),
		"user_id":         appcontext.GetUserID(ctx),
	}

	if len(fields) > 0 {
		logFields["data"] = fields[0]
	} else {
		logFields["data"] = nil
	}

	return logFields
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
	} else {
		location := strings.Split(file, "/")
		file = strings.Join(location[len(location)-3:], "/")
	}
	return fmt.Sprintf("%s:%d", file, line)
}
