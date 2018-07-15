package base

import (
	"github.com/go-sql-driver/mysql"
)

func IsMySQLDuplicateEntryError(err error) bool {
	const (
		mysqlDuplicateEntryError = 1062
	)

	switch err.(type) {
	case *mysql.MySQLError:
		mysqlerr, ok := err.(*mysql.MySQLError)
		if ok && mysqlerr.Number == mysqlDuplicateEntryError {
			return true
		}
	}

	return false
}
