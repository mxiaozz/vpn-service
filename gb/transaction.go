package gb

import "xorm.io/xorm"

func Tx(action func(dbSession *xorm.Session) error) error {
	dbSession := DB.NewSession()
	defer dbSession.Close()
	if err := dbSession.Begin(); err != nil {
		return err
	}
	if err := action(dbSession); err != nil {
		return err
	}
	return dbSession.Commit()
}

func TxWith[T any](action func(dbSession *xorm.Session) (T, error)) (T, error) {
	var zero T
	dbSession := DB.NewSession()
	defer dbSession.Close()
	if err := dbSession.Begin(); err != nil {
		return zero, err
	}
	if result, err := action(dbSession); err != nil {
		return result, err
	}
	return zero, dbSession.Commit()
}
