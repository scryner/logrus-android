package logrus_android

/*
#cgo LDFLAGS: -landroid -llog

#include <android/log.h>
#include <string.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/Sirupsen/logrus"
)

var allLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
}

type androidHook struct {
	tag *C.char
	fmt logrus.Formatter
}

func (ah *androidHook) Levels() []logrus.Level {
	return allLevels
}

func (ah *androidHook) Fire(e *logrus.Entry) error {
	var priority C.int

	formatted, err := ah.fmt.Format(e)
	if err != nil {
		return err
	}
	cstr := C.CString(string(formatted))

	switch e.Level {
	case logrus.PanicLevel:
		priority = C.ANDROID_LOG_FATAL
	case logrus.FatalLevel:
		priority = C.ANDROID_LOG_FATAL
	case logrus.ErrorLevel:
		priority = C.ANDROID_LOG_ERROR
	case logrus.WarnLevel:
		priority = C.ANDROID_LOG_WARN
	case logrus.InfoLevel:
		priority = C.ANDROID_LOG_INFO
	case logrus.DebugLevel:
		priority = C.ANDROID_LOG_DEBUG
	}
	C.__android_log_write(priority, ah.tag, cstr)
	C.free(unsafe.Pointer(cstr))
	return nil
}

// NewAndroidHook create a logrus Hook that forward entries to logcat
func NewAndroidHook(tag string) logrus.Hook {
	return &androidHook{
		tag: C.CString(tag),
		fmt: &logrus.TextFormatter{
			ForceColors:      false,
			DisableColors:    true,
			DisableTimestamp: true,
			FullTimestamp:    false,
			DisableSorting:   false,
		},
	}
}
