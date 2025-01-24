package common

import "github.com/charmbracelet/log"

func Log(msg interface{}, keyvals ...interface{}) {
	log.Info(msg, keyvals...)
}

func Logf(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func LogError(err error, keyvals ...interface{}) {
	log.Error(err, keyvals...)
}

func LogErrorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func LogWarn(msg interface{}, keyvals ...interface{}) {
	log.Warn(msg, keyvals...)
}

func LogWarnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func LogDebug(msg interface{}, keyvals ...interface{}) {
	log.Debug(msg, keyvals...)
}

func LogDebugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}
