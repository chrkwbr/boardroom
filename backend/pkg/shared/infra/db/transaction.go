package db

import "database/sql"

type Transaction interface {
	RunWithTx(fn func(tx *sql.Tx) error) error
}

type TransactionManagerImpl struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) Transaction {
	return &TransactionManagerImpl{db: db}
}

func (tm *TransactionManagerImpl) RunWithTx(fn func(tx *sql.Tx) error) error {
	tx, err := tm.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				panic(rollbackErr)
			}
			panic(r)
		}
	}()

	if err = fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit()
}
