package usecase

import (
	"fmt"
	"sync"
	"time"
)

const (
	putSessionFileRateWindow       = time.Minute
	putSessionFileMaxBytesPerMin   = 48 * 1024 * 1024
	putSessionFileMaxUploadsPerMin = 80
	putSessionFileRatePruneEntries = 4096
)

type sessionPutRateLimiter struct {
	mu      sync.Mutex
	perUser map[int]*uploadRollingWindow
	now     func() time.Time
}

type uploadRollingWindow struct {
	windowStart time.Time
	bytes       int64
	n           int
}

func (l *sessionPutRateLimiter) currentTime() time.Time {
	if l.now != nil {
		return l.now()
	}

	return time.Now()
}

func (l *sessionPutRateLimiter) checkPutSessionFileRate(userID int, byteLen int) error {
	if byteLen < 0 {
		byteLen = 0
	}

	b := int64(byteLen)
	now := l.currentTime()

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.perUser == nil {
		l.perUser = make(map[int]*uploadRollingWindow)
	}

	w := l.perUser[userID]
	if w == nil || now.Sub(w.windowStart) >= putSessionFileRateWindow {
		l.perUser[userID] = &uploadRollingWindow{
			windowStart: now,
			bytes:       b,
			n:           1,
		}

		l.maybePruneLocked(now)

		return nil
	}

	if w.n >= putSessionFileMaxUploadsPerMin {
		return fmt.Errorf("слишком много загрузок файлов за минуту (лимит %d)", putSessionFileMaxUploadsPerMin)
	}

	if w.bytes+b > putSessionFileMaxBytesPerMin {
		return fmt.Errorf(
			"превышен объём загрузок за минуту: лимит %d МиБ",
			putSessionFileMaxBytesPerMin/(1024*1024),
		)
	}

	w.bytes += b
	w.n++
	l.maybePruneLocked(now)
	return nil
}

func (l *sessionPutRateLimiter) maybePruneLocked(now time.Time) {
	if len(l.perUser) < putSessionFileRatePruneEntries {
		return
	}

	for uid, w := range l.perUser {
		if w == nil || now.Sub(w.windowStart) >= putSessionFileRateWindow {
			delete(l.perUser, uid)
		}
	}
}
