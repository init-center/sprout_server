package transaction

import (
	"github.com/jmoiron/sqlx"
)

type TxFunc func(*sqlx.Tx) (err error)

func Start(db *sqlx.DB, txFunc TxFunc) (err error) {
	// start a transaction
	tx, err := db.Beginx()
	if err != nil {
		return
	}

	// if return, always rollback the transaction
	// but if committed, cannot be rolled back
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// run txFunc
	err = txFunc(tx)

	return
}
