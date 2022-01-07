package logger

import (
    "context"
    "fmt"
    "time"
)

func NewSimpleJSONLogger(appName, appInstanceID string) Logger {
	return &simpleJSONLoggerImpl{
		ApplicationData: ApplicationData{
			AppName:       appName,
			AppInstanceID: appInstanceID,
			StartTime:     time.Now().Format("2006-01-02 15:04:05"),
		},
	}
}

type jsonLogModel struct {
	AppName   string `json:"appName"`
	AppInstID string `json:"appInstID"`
	Start     string `json:"start"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	Location  string `json:"location"`
	Time      string `json:"time"`
}

func newJSONLogModel(lg *simpleJSONLoggerImpl, flag, loc string, msg, trid interface{}) string {

	if flag == "ERROR" {
		return toJsonString(jsonLogModel{
			AppName:   lg.AppName,
			AppInstID: lg.AppInstanceID,
			Start:     lg.StartTime,
			Severity:  flag,
			Message:   fmt.Sprintf("%v %v %v", trid, loc, msg),
			Location:  loc,
			Time:      time.Now().String(),
		})
	}

	return toJsonString(jsonLogModel{
		AppName:   lg.AppName,
		AppInstID: lg.AppInstanceID,
		Start:     lg.StartTime,
		Severity:  flag,
		Message:   fmt.Sprintf("%v %v", trid, msg),
		Location:  loc,
		Time:      time.Now().String(),
	})
}

type simpleJSONLoggerImpl struct {
	ApplicationData
}

func (l simpleJSONLoggerImpl) Info(ctx context.Context, message string, args ...interface{}) {
	messageWithArgs := fmt.Sprintf(message, args...)
	l.printLog(ctx, "INFO", messageWithArgs)
}

func (l simpleJSONLoggerImpl) Error(ctx context.Context, message string, args ...interface{}) {
	messageWithArgs := fmt.Sprintf(message, args...)
	l.printLog(ctx, "ERROR", messageWithArgs)
}

func (l simpleJSONLoggerImpl) GetApplicationData() ApplicationData {
	return l.ApplicationData
}

func (l simpleJSONLoggerImpl) printLog(ctx context.Context, flag string, data interface{}) {

	// default traceId
	traceID := "0000000000000000"

	if ctx != nil {
		if v := ctx.Value(traceDataKey); v != nil {
			traceID = v.(string)
		}
	}

	fmt.Println(newJSONLogModel(&l, flag, getFileLocationInfo(3), data, traceID))

}
