package db

import (
	"database/sql"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// TrxSupportRepo database trx support
type TrxSupportRepo interface {
	Begin() (TrxObj, error)
}

type TrxObj interface {
	Commit() error
	Rollback() error
}

// GormTrxSupport parent mysqlrepo
type GormTrxSupport struct {
	DB *gorm.DB
}

type GormTrxObj struct {
	DB *gorm.DB
}

// Begin Begin db transaction
func (repo *GormTrxSupport) Begin() (TrxObj, error) {
	trx := repo.DB.Begin(&sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})

	return &GormTrxObj{DB: trx}, trx.Error
}

// Trx get transaction
func (repo *GormTrxSupport) Trx(trx TrxObj) *gorm.DB {
	gormTrx, ok := trx.(*GormTrxObj)
	if ok {
		return gormTrx.DB
	}

	return repo.DB
}

// Commit Commit db transaction
func (trx *GormTrxObj) Commit() error {
	return trx.DB.Commit().Error
}

// Rollback rollback trx
func (trx *GormTrxObj) Rollback() error {
	return trx.DB.Rollback().Error
}

// DBTransaction usecase with db transaction
func DBTransaction(repo TrxSupportRepo, callback func(TrxObj) error) (err error) {
	functionName := "DBTransaction"
	commit := false
	trx, err := repo.Begin()
	if err != nil {
		return err
	}

	defer func(commit *bool, repo TrxSupportRepo, trx TrxObj) {
		if !*commit {
			if rErr := trx.Rollback(); rErr != nil {
				if err == nil {
					err = rErr
				} else {
					err = errors.Wrap(rErr, err.Error())
				}
			}
		}
	}(&commit, repo, trx)

	// call the callback function
	err = callback(trx)
	if err != nil {
		return err
	}

	if err = trx.Commit(); err != nil {
		return errors.Wrap(err, functionName)
	}
	commit = true

	return err
}
