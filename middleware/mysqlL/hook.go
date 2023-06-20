package mysqlL

import (
	"context"
	"github.com/flyerxp/lib/logger"
	"go.uber.org/zap"
	"strconv"
	"sync/atomic"
	"time"
)

type Hooks struct {
	*zap.Logger
	sqlKey             *sqlDurationKey
	IsPrintSQLDuration bool
}
type sqlDurationKey struct {
	nums int32
}

func (s *sqlDurationKey) String() string {
	return "query_" + strconv.Itoa(int(s.nums))
}
func (s *sqlDurationKey) add() {
	atomic.AddInt32(&s.nums, 1)
}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if h.sqlKey == nil {
		h.sqlKey = &sqlDurationKey{0}
	}
	h.sqlKey.add()
	if h.IsPrintSQLDuration {
		logger.AddNotice(zap.String(h.sqlKey.String(), query), zap.Any(h.sqlKey.String(), args))
	}
	return context.WithValue(ctx, h.sqlKey.String(), time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin, ok := ctx.Value(h.sqlKey.String()).(time.Time)
	if ok {
		timeout := int(time.Since(begin).Milliseconds())
		logger.AddMysqlTime(timeout)
		if timeout > 2000 {
			logger.AddWarn(zap.String(h.sqlKey.String(), query), zap.Any(h.sqlKey.String(), args))
		}
	}
	return ctx, nil
}
func (h *Hooks) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	if begin, ok := ctx.Value(h.sqlKey.String()).(time.Time); ok {
		logger.AddMysqlTime(int(time.Since(begin).Milliseconds()))
	}
	logger.AddError(zap.String(h.sqlKey.String(), query), zap.Any(h.sqlKey.String(), args), zap.Error(err))
	return nil
}
