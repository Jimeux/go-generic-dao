package test

import (
	"testing"

	"github.com/Jimeux/go-generic-dao/db"
)

func Truncate(t *testing.T, table string) func() {
	return func() {
		if _, err := db.DB().Exec("TRUNCATE TABLE " + table + ";"); err != nil {
			t.Logf("failed to truncate table %s: %v", table, err)
		}
	}
}
