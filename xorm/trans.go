package xorm

import (
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/jtolds/gls"
	"sync"
)

func newTrans(logger gone.Logger, newSession func() Session) trans {
	return trans{
		logger:     logger,
		newSession: newSession,
	}
}

type trans struct {
	logger     gone.Logger
	newSession func() Session
}

var sessionMap = sync.Map{}

// ============================================================================
func (e *trans) getTransaction(id uint) (Session, bool) {
	session, suc := sessionMap.Load(id)
	if suc {
		return session.(Session), false
	} else {
		s := e.newSession()
		sessionMap.Store(id, s)
		return s, true
	}
}

func (e *trans) delTransaction(id uint, session Session) error {
	defer sessionMap.Delete(id)
	return session.Close()
}

// Transaction execute sql in transaction
func (e *trans) Transaction(fn func(session Interface) error) error {
	var err error
	gls.EnsureGoroutineId(func(gid uint) {
		session, isNew := e.getTransaction(gid)

		if isNew {
			rollback := func() {
				rollbackErr := session.Rollback()
				if rollbackErr != nil {
					e.logger.Errorf("rollback err:%v", rollbackErr)
					err = gone.ToErrorWithMsg(err, fmt.Sprintf("rollback error: %v", rollbackErr))
				}
			}

			isRollback := false
			defer func(e *trans, id uint, session Session) {
				err := e.delTransaction(id, session)
				if err != nil {
					e.logger.Errorf("del session err:%v", err)
				}
			}(e, gid, session)

			defer func() {
				if info := recover(); info != nil {
					e.logger.Errorf("session rollback for panic: %s", info)
					e.logger.Errorf("%s", gone.PanicTrace(2, 1))
					if !isRollback {
						rollback()
						err = gone.NewInnerError(fmt.Sprintf("%s", info), gone.DbRollForPanicError)
					} else {
						err = gone.ToErrorWithMsg(info, fmt.Sprintf("rollback for err: %v, but panic for", err))
					}
				}
			}()

			err = session.Begin()
			if err != nil {
				err = gone.ToError(err)
				return
			}
			err = gone.ToError(fn(session))
			if err == nil {
				err = gone.ToError(session.Commit())
			} else {
				e.logger.Errorf("session rollback for err: %v", err)
				isRollback = true
				rollback()
			}
		} else {
			err = gone.ToError(fn(session))
		}
	})
	return err
}
