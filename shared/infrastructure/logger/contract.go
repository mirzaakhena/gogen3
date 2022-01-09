package logger

import (
    "context"
    "encoding/json"
    "fmt"
    "runtime"
    "strings"
)

type ApplicationData struct {
	AppName       string `json:"appName"`
	AppInstanceID string `json:"appInstanceID"`
	StartTime     string `json:"startTime"`
}

type Logger interface {
	Info(ctx context.Context, message string, args ...interface{})
	Error(ctx context.Context, message string, args ...interface{})
	GetApplicationData() ApplicationData
}

type traceDataType int

const traceDataKey traceDataType = 1

func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceDataKey, traceID)
}

// getFileLocationInfo get the function information like filename and line number
// skip is the parameter that need to adjust if we add new method layer
func getFileLocationInfo(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	funcName := runtime.FuncForPC(pc).Name()
	x := strings.LastIndex(funcName, "/")
	return fmt.Sprintf("%s:%d", funcName[x+1:], line)
}

func toJsonString(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}