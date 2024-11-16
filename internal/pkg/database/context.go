package database

import (
	"github.com/jmoiron/sqlx"
	"sync"
)

type TransactionCtx struct {
	Mu   sync.Mutex
	Conn *sqlx.Tx
}

func (t *TransactionCtx) Commit() error {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if t.Conn != nil {
		return t.Conn.Commit()
	}
	return nil
}

func (t *TransactionCtx) Rollback() error {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if t.Conn != nil {
		return t.Conn.Rollback()
	}
	return nil
}

type CustomSettingCtx struct {
	IsJobAfterTxCommit bool
}

type TransactionCtxKeyType string

type CustomSettingCtxKeyType string

const (
	TransactionCtxKey TransactionCtxKeyType = "transactionCtx"
)

const (
	CustomSettingCtxKey CustomSettingCtxKeyType = "customSettingCtx"
)
