package logger

import (
	"go.uber.org/zap"
)

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

/*
	func AddCtxNotice(ctx context.Context, field ...zap.Field) {
		if !noticeLog.isInitEd {
			getNoticeLog()
		}
		//field = append(field, zap.String(ctx.Value("requestId")))
		noticeLog.noticeMetrics.Notice = append(noticeLog.noticeMetrics.Notice, field...)
	}

	func AddCtxError(ctx context.Context, field ...zap.Field) {
		if !errLogV.isInitEd {
			getErrorLog()
		}
		errLogV.errMetrics.Error = append(errLogV.errMetrics.Error, field...)
	}

	func AddCtxWarn(ctx context.Context, field ...zap.Field) {
		if !warnLogV.isInitEd {
			getWarnLog()
		}
		warnLogV.warnMetrics.Warn = append(warnLogV.warnMetrics.Warn, field...)
	}
*/
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
