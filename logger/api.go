package logger

import (
	"go.uber.org/zap"
)

func AddNotice(field ...zap.Field) {
	if !noticeLog.isInitEd {
		GetNoticeLog()
	}
	//field = append(field, zap.String(ctx.Value("requestId")))
	noticeLog.noticeMetrics.Notice = append(noticeLog.noticeMetrics.Notice, field...)
}
func AddError(field ...zap.Field) {
	if !errLogV.isInitEd {
		GetErrorLog()
	}
	errLogV.errMetrics.Error = append(errLogV.errMetrics.Error, field...)
}
func AddWarn(field ...zap.Field) {
	if !warnLogV.isInitEd {
		GetWarnLog()
	}
	warnLogV.warnMetrics.Warn = append(warnLogV.warnMetrics.Warn, field...)
}

/*
	func AddCtxNotice(ctx context.Context, field ...zap.Field) {
		if !noticeLog.isInitEd {
			GetNoticeLog()
		}
		//field = append(field, zap.String(ctx.Value("requestId")))
		noticeLog.noticeMetrics.Notice = append(noticeLog.noticeMetrics.Notice, field...)
	}

	func AddCtxError(ctx context.Context, field ...zap.Field) {
		if !errLogV.isInitEd {
			GetErrorLog()
		}
		errLogV.errMetrics.Error = append(errLogV.errMetrics.Error, field...)
	}

	func AddCtxWarn(ctx context.Context, field ...zap.Field) {
		if !warnLogV.isInitEd {
			GetWarnLog()
		}
		warnLogV.warnMetrics.Warn = append(warnLogV.warnMetrics.Warn, field...)
	}
*/
func WriteLine() {
	noticeLog.ZapLog.With(zap.Int("execTime", noticeLog.noticeMetrics.TotalExecTime)).With(zap.Object("middle", noticeLog.noticeMetrics.Middle)).With(noticeLog.noticeMetrics.Notice...).Info("info")
	if len(errLogV.errMetrics.Error) > 1 {
		errLogV.ZapLog.With(zap.Int("execTime", noticeLog.noticeMetrics.TotalExecTime)).With(zap.Object("middle", noticeLog.noticeMetrics.Middle)).With(errLogV.errMetrics.Error...).WithOptions(zap.AddCallerSkip(1)).Error("error")
	}
	if len(warnLogV.warnMetrics.Warn) > 1 {
		warnLogV.ZapLog.With(zap.Int("execTime", noticeLog.noticeMetrics.TotalExecTime)).With(warnLogV.warnMetrics.Warn...).WithOptions(zap.AddCallerSkip(1)).Warn("warn")
	}
	Reset()
}

func Reset() {
	noticeLog.noticeMetrics.Middle = MiddleExecTime{}
	noticeLog.noticeMetrics.Notice = make([]zap.Field, 1, 10)
	noticeLog.noticeMetrics.Notice[0] = zap.Namespace("notice")
	noticeLog.noticeMetrics.TotalExecTime = 0
	errLogV.errMetrics.Error = make([]zap.Field, 1, 10)
	errLogV.errMetrics.Error[0] = zap.Namespace("error")
	warnLogV.warnMetrics.Warn = make([]zap.Field, 1, 10)
	warnLogV.warnMetrics.Warn[0] = zap.Namespace("warn")
}
