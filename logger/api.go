package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func init() {
	go func() {
		getNoticeLog()
		getErrorLog()
		getWarnLog()
	}()
}
func AddNotice(field ...zap.Field) {
	if !noticeLog.isInitEd {
		getNoticeLog()
	}
	//field = append(field, zap.String(ctx.Value("requestId")))
	noticeLog.noticeMetrics.Notice = append(noticeLog.noticeMetrics.Notice, field...)
}
func AddError(field ...zap.Field) {
	if !errLogV.isInitEd {
		getErrorLog()
	}
	/*for i, _ := range field {
		warnLogV.ZapLog.Error("error", field[i])
	}*/
	errLogV.errMetrics.Error = append(errLogV.errMetrics.Error, field...)
}
func AddWarn(field ...zap.Field) {
	if !warnLogV.isInitEd {
		getWarnLog()
	}
	/*for i, _ := range field {
		warnLogV.ZapLog.Warn("warn", field[i])
	}*/
	warnLogV.warnMetrics.Warn = append(warnLogV.warnMetrics.Warn, field...)
}

var CurrDay = time.Now().Day()

type Event struct {
	Name string
	F    func()
}

var makeFileEvent = make([]Event, 0, 15)

var resetEvent = make([]Event, 0, 15)

// 生成文件的时候执行
func RegistermakeFileEvent(F Event) {
	makeFileEvent = append(makeFileEvent, F)
}

// Reset 的时候执行
func RegisterReset(F Event) {
	resetEvent = append(resetEvent, F)
}
func GetPath(paths []string, vtype string) []string {
	pathNew := make([]string, 0, len(paths))
	//copy(pathNew, paths)
	for i := range paths {
		if (paths[i] == "stdout" || paths[i] == "stderr") && runtime.GOOS == "windows" {
			continue
		}
		if paths[i] != "" && paths[i] != "stdout" && paths[i] != "stderr" && strings.Contains(paths[i], "/") {
			pathTmp := filepath.Dir(paths[i])
			_, err := os.Stat(pathTmp)
			if os.IsNotExist(err) {
				_ = os.MkdirAll(pathTmp, os.ModePerm)
			}
			fmt.Println("怎么回事")
			pathNew = append(pathNew, paths[i]+fmt.Sprintf("_%s_%02d_%d.log", vtype, time.Now().Month(), time.Now().Day()))
		}
	}
	fmt.Println(pathNew, paths)
	return pathNew
}
func WriteLine() {
	if !noticeLog.isInitEd {
		getNoticeLog()
	}
	noticeLog.ZapLog.With(zap.Int("execTotalTime", noticeLog.noticeMetrics.TotalExecTime)).With(zap.Object("middle", noticeLog.noticeMetrics.Middle)).With(zap.Object("execTime", noticeLog.execMetrics)).With(noticeLog.noticeMetrics.Notice...).Info("info")
	if len(errLogV.errMetrics.Error) > 1 {
		if !errLogV.isInitEd {
			getErrorLog()
		}
		errLogV.ZapLog.With(zap.Int("execTotalTime", noticeLog.noticeMetrics.TotalExecTime)).With(zap.Object("middle", noticeLog.noticeMetrics.Middle)).With(errLogV.errMetrics.Error...).WithOptions(zap.AddCallerSkip(1)).Error("error")
	}
	if len(warnLogV.warnMetrics.Warn) > 1 {
		if !warnLogV.isInitEd {
			getWarnLog()
		}
		warnLogV.ZapLog.With(zap.Int("execTotalTime", noticeLog.noticeMetrics.TotalExecTime)).With(warnLogV.warnMetrics.Warn...).WithOptions(zap.AddCallerSkip(1)).Warn("warn")
	}
	//第二天重置
	if CurrDay != time.Now().Day() {
		makeFileEventNew := make([]Event, 0, 15)
		copy(makeFileEventNew, makeFileEvent)
		makeFileEvent = make([]Event, 0, 15)
		resetEvent = make([]Event, 0, 15)
		Reset()
		for _, f := range makeFileEventNew {
			f.F()
		}
		CurrDay = time.Now().Day()
	}
	Reset()
}
func WriteErr() {
	if len(errLogV.errMetrics.Error) > 1 {
		if !errLogV.isInitEd {
			getErrorLog()
		}
		errLogV.ZapLog.With(zap.Int("execTotalTime", noticeLog.noticeMetrics.TotalExecTime)).With(zap.Object("middle", noticeLog.noticeMetrics.Middle)).With(errLogV.errMetrics.Error...).WithOptions(zap.AddCallerSkip(1)).Error("error")
		errLogV.errMetrics.Error = make([]zap.Field, 1, 10)
		errLogV.errMetrics.Error[0] = zap.Namespace("error")
	}
	if len(warnLogV.warnMetrics.Warn) > 1 {
		if !warnLogV.isInitEd {
			getWarnLog()
		}
		warnLogV.ZapLog.With(zap.Int("execTotalTime", noticeLog.noticeMetrics.TotalExecTime)).With(warnLogV.warnMetrics.Warn...).WithOptions(zap.AddCallerSkip(1)).Warn("warn")
		warnLogV.warnMetrics.Warn = make([]zap.Field, 1, 10)
		warnLogV.warnMetrics.Warn[0] = zap.Namespace("warn")
	}
	Reset()
}

func Reset() {
	for _, f := range resetEvent {
		f.F()
	}
}
