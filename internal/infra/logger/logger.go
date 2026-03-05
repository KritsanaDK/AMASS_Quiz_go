package logger

import (
	"amass/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerCommonInstance ILogSystem
)

type (
	LoggerLevel uint32

	ILogSystem interface {
		Info(message string, args ...interface{})
		Warn(message string, args ...interface{})
		Debug(message string, args ...interface{})
		Error(message string, args ...interface{})
		Trace(message string, args ...interface{})

		InfoWithFields(message string, fields ...zap.Field)
		WarnWithFields(message string, fields ...zap.Field)
		DebugWithFields(message string, fields ...zap.Field)
		ErrorWithFields(message string, fields ...zap.Field)
		TraceWithFields(message string, fields ...zap.Field)

		InfoWithTracking(message string, tracking *models.Tracking)
		WarnWithTracking(message string, tracking *models.Tracking)
		DebugWithTracking(message string, tracking *models.Tracking)
		ErrorWithTracking(message, errMessage string, tracking *models.Tracking)
		TraceWithTracking(message string, tracking *models.Tracking)

		Sync()
		GetSugar() *zap.SugaredLogger
	}

	logSystem struct {
		loggerInstance        *zap.Logger
		sugaredLoggerInstance *zap.SugaredLogger
	}

	LogSystemType interface {
	}
)

func InitialzeLoggerSystem() error {
	if loggerCommonInstance == nil {
		// initialize log
		configLog := zap.NewProductionConfig()
		configLog.DisableCaller = true
		configLog.DisableStacktrace = true
		configLog.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

		zapLog, err := configLog.Build()
		if err != nil {
			return fmt.Errorf("CreateLogger(%+v): %v", configLog, err)
		}
		sugar := zapLog.Sugar()

		loggerCommonInstance = &logSystem{
			loggerInstance:        zapLog,
			sugaredLoggerInstance: sugar,
		}
	}

	return nil
}

func Sync() {
	if loggerCommonInstance != nil {
		loggerCommonInstance.Sync()
	}
}

func GetSugar() *zap.SugaredLogger {
	if loggerCommonInstance != nil {
		return loggerCommonInstance.GetSugar()
	}

	return nil
}

func Info(message string, args ...interface{}) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.Info(message, args...)
	} else {
		fmt.Printf("[Info] "+message, args...)
	}
}
func Warn(message string, args ...interface{}) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.Warn(message, args...)
	} else {
		fmt.Printf("[Warn] "+message, args...)
	}
}
func Debug(message string, args ...interface{}) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.Debug(message, args...)
	} else {
		fmt.Printf("[Debug] "+message, args...)
	}
}
func Error(message string, args ...interface{}) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.Error(message, args...)
	} else {
		fmt.Printf("[Error] "+message, args...)
	}
}
func Trace(message string, args ...interface{}) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.Trace(message, args...)
	} else {
		fmt.Printf("[Trace] "+message, args...)
	}
}

func InfoWithFields(message string, fields ...zap.Field) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.InfoWithFields(message, fields...)
	} else {
		fieldsString := ""
		prettyJSON, err := json.MarshalIndent(fields, "", "  ")
		if err == nil {
			fieldsString = string(prettyJSON)
		}

		fmt.Printf("[Into] "+message+" %+v", fieldsString)
	}
}
func WarnWithFields(message string, fields ...zap.Field) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.WarnWithFields(message, fields...)
	} else {
		fieldsString := ""
		prettyJSON, err := json.MarshalIndent(fields, "", "  ")
		if err == nil {
			fieldsString = string(prettyJSON)
		}

		fmt.Printf("[Warn] "+message+" %+v", fieldsString)
	}
}
func DebugWithFields(message string, fields ...zap.Field) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.DebugWithFields(message, fields...)
	} else {
		fieldsString := ""
		prettyJSON, err := json.MarshalIndent(fields, "", "  ")
		if err == nil {
			fieldsString = string(prettyJSON)
		}

		fmt.Printf("[Debug] "+message+" %+v", fieldsString)
	}
}
func ErrorWithFields(message string, fields ...zap.Field) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.ErrorWithFields(message, fields...)
	} else {
		fieldsString := ""
		prettyJSON, err := json.MarshalIndent(fields, "", "  ")
		if err == nil {
			fieldsString = string(prettyJSON)
		}

		fmt.Printf("[Error] "+message+" %+v", fieldsString)
	}
}
func TraceWithFields(message string, fields ...zap.Field) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.TraceWithFields(message, fields...)
	} else {
		fieldsString := ""
		prettyJSON, err := json.MarshalIndent(fields, "", "  ")
		if err == nil {
			fieldsString = string(prettyJSON)
		}

		fmt.Printf("[Trace] "+message+" %+v", fieldsString)
	}
}

func InfoWithTracking(message string, tracking *models.Tracking) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.InfoWithTracking(message, tracking)
	} else {
		fmt.Printf("[Info] "+message+" %+v", tracking)
	}
}
func WarnWithTracking(message string, tracking *models.Tracking) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.WarnWithTracking(message, tracking)
	} else {
		fmt.Printf("[Warn] "+message+" %+v", tracking)
	}
}
func DebugWithTracking(message string, tracking *models.Tracking) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.DebugWithTracking(message, tracking)
	} else {
		fmt.Printf("[Debug] "+message+" %+v", tracking)
	}
}
func ErrorWithTracking(message, errMessage string, tracking *models.Tracking) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.ErrorWithTracking(message, errMessage, tracking)
	} else {
		fmt.Printf("[Error] "+message+" (%s) %+v", errMessage, tracking)
	}
}
func TraceWithTracking(message string, tracking *models.Tracking) {
	if loggerCommonInstance != nil {
		loggerCommonInstance.TraceWithTracking(message, tracking)
	} else {
		fmt.Printf("[Trace] "+message+" %+v", tracking)
	}
}

////////////////////////////////////////////////////////////////
// instance function
////////////////////////////////////////////////////////////////

func (log *logSystem) Sync() {
	log.loggerInstance.Sync()
}

func (log *logSystem) GetSugar() *zap.SugaredLogger {
	return log.sugaredLoggerInstance
}

func (log *logSystem) Info(message string, args ...interface{}) {
	log.sugaredLoggerInstance.Infof(message, args...)
}
func (log *logSystem) Debug(message string, args ...interface{}) {
	log.sugaredLoggerInstance.Debugf(message, args...)
}
func (log *logSystem) Error(message string, args ...interface{}) {
	log.sugaredLoggerInstance.Errorf(message, args...)
}
func (log *logSystem) Warn(message string, args ...interface{}) {
	log.sugaredLoggerInstance.Warnf(message, args...)
}
func (log *logSystem) Trace(message string, args ...interface{}) {
	log.sugaredLoggerInstance.Debugf(message, args...)
}

func (log *logSystem) InfoWithFields(message string, fields ...zap.Field) {
	log.loggerInstance.Info(message, fields...)
}
func (log *logSystem) DebugWithFields(message string, fields ...zap.Field) {
	log.loggerInstance.Debug(message, fields...)
}
func (log *logSystem) ErrorWithFields(message string, fields ...zap.Field) {
	log.loggerInstance.Error(message, fields...)
}
func (log *logSystem) WarnWithFields(message string, fields ...zap.Field) {
	log.loggerInstance.Warn(message, fields...)
}
func (log *logSystem) TraceWithFields(message string, fields ...zap.Field) {
	log.loggerInstance.Debug(message, fields...)
}

func (log *logSystem) InfoWithTracking(message string, tracking *models.Tracking) {
	now := time.Now()
	dateTime := now.Format(time.DateTime)
	interval := now.Sub(tracking.StartDate).Seconds()
	//message specific handler or service
	log.sugaredLoggerInstance.Infow(message,
		zap.String("uri", tracking.URI),
		zap.String("method", tracking.Method),
		zap.Any("request", tracking.Request),
		zap.Any("response", tracking.Response),
		zap.String("date_time", dateTime),
		zap.String("track_id", tracking.Track),
		zap.String("interval", fmt.Sprintf("%f", interval)),
	)
}
func (log *logSystem) WarnWithTracking(message string, tracking *models.Tracking) {
	now := time.Now()
	dateTime := now.Format(time.DateTime)
	interval := now.Sub(tracking.StartDate).Seconds()
	//message specific handler or service
	log.sugaredLoggerInstance.Warnw(message,
		zap.String("uri", tracking.URI),
		zap.String("method", tracking.Method),
		zap.Any("request", tracking.Request),
		zap.Any("response", tracking.Response),
		zap.String("date_time", dateTime),
		zap.String("track_id", tracking.Track),
		zap.String("interval", fmt.Sprintf("%f", interval)),
	)
}
func (log *logSystem) DebugWithTracking(message string, tracking *models.Tracking) {
	now := time.Now()
	dateTime := now.Format(time.DateTime)
	interval := now.Sub(tracking.StartDate).Seconds()
	//message specific handler or service
	log.sugaredLoggerInstance.Debug(message,
		zap.String("uri", tracking.URI),
		zap.String("method", tracking.Method),
		zap.Any("request", tracking.Request),
		zap.Any("response", tracking.Response),
		zap.String("date_time", dateTime),
		zap.String("track_id", tracking.Track),
		zap.String("interval", fmt.Sprintf("%f", interval)),
	)
}
func (log *logSystem) ErrorWithTracking(message, errMessage string, tracking *models.Tracking) {
	now := time.Now()
	dateTime := now.Format(time.DateTime)
	interval := now.Sub(tracking.StartDate).Seconds()
	//message specific handler or service
	log.sugaredLoggerInstance.Errorw(message,
		zap.String("uri", tracking.URI),
		zap.String("method", tracking.Method),
		zap.Any("request", tracking.Request),
		zap.Any("response", tracking.Response),
		zap.String("message", errMessage),
		zap.String("date_time", dateTime),
		zap.String("track_id", tracking.Track),
		zap.String("interval", fmt.Sprintf("%f", interval)),
	)
}
func (log *logSystem) TraceWithTracking(message string, tracking *models.Tracking) {
	now := time.Now()
	dateTime := now.Format(time.DateTime)
	interval := now.Sub(tracking.StartDate).Seconds()
	//message specific handler or service
	log.sugaredLoggerInstance.Debug(message,
		zap.String("uri", tracking.URI),
		zap.String("method", tracking.Method),
		zap.Any("request", tracking.Request),
		zap.Any("response", tracking.Response),
		zap.String("date_time", dateTime),
		zap.String("track_id", tracking.Track),
		zap.String("interval", fmt.Sprintf("%f", interval)),
	)
}

func LogInfoQueue(message, infoMessage string, log *zap.SugaredLogger, tracking *models.Tracking) {
	now := time.Now()
	dateTime := now.Format(time.DateTime)
	interval := now.Sub(tracking.StartDate).Seconds()
	log.Infow(message,
		zap.Any("request", tracking.Request),
		zap.String("message", infoMessage),
		zap.String("date_time", dateTime),
		zap.Int64("unix", time.Now().UnixMicro()),
		zap.String("interval", fmt.Sprintf("%f", interval)),
	)
}

func LogErrorQueue(message, errMessage string, log *zap.SugaredLogger, tracking *models.Tracking) {
	now := time.Now()
	dateTime := now.Format(time.DateTime)
	interval := now.Sub(tracking.StartDate).Seconds()
	log.Errorw(message,
		zap.Any("request", tracking.Request),
		zap.String("message", errMessage),
		zap.String("date_time", dateTime),
		zap.String("interval", fmt.Sprintf("%f", interval)),
	)
}
