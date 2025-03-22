package g

import "time"

type DoLocker interface {
	LockAndDo(key string, fn func(), lockTime, checkPeriod time.Duration) (err error)
}
