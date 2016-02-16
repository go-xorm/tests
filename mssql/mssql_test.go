package tests

import (
	"database/sql"
	"runtime"
	"testing"

	. ".."
	"github.com/go-xorm/xorm"
	_ "github.com/lunny/godbc"
)

var mssqlConnStr string

func init() {
	if runtime.GOOS == "windows" {
		mssqlConnStr = "driver={SQL Server};Server=192.168.3.103;Database=xorm_test; uid=sa; pwd=1234;"
	} else {
		mssqlConnStr = "driver={freetds};Server=192.168.3.103;Database=xorm_test; uid=sa; pwd=1234;"
	}
}

func newMssqlEngine() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("odbc", mssqlConnStr)
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(ShowTestSql)
	return engine, nil
}

func TestMssql(t *testing.T) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	BaseTestAll(engine, t)
	BaseTestAll2(engine, t)
}

func TestMssqlWithCache(t *testing.T) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())

	BaseTestAll(engine, t)
	BaseTestAll2(engine, t)
}

func newMssqlDriverDB() (*sql.DB, error) {
	return sql.Open("odbc", mssqlConnStr)
}

const (
	createTableMssql = `IF NOT EXISTS (SELECT [name] FROM sys.tables WHERE [name] = 'big_struct' ) CREATE TABLE
		"big_struct" ("id" BIGINT PRIMARY KEY IDENTITY NOT NULL, "name" VARCHAR(255) NULL, "title" VARCHAR(255) NULL,
		"age" VARCHAR(255) NULL, "alias" VARCHAR(255) NULL, "nick_name" VARCHAR(255) NULL);
		`

	dropTableMssql = "IF EXISTS (SELECT * FROM sysobjects WHERE id = object_id(N'big_struct') and OBJECTPROPERTY(id, N'IsUserTable') = 1) DROP TABLE IF EXISTS `big_struct`;"
)

func BenchmarkMssqlDriverInsert(t *testing.B) {
	DoBenchDriver(newMssqlDriverDB, createTableMssql, dropTableMssql,
		DoBenchDriverInsert, t)
}

func BenchmarkMssqlDriverFind(t *testing.B) {
	DoBenchDriver(newMssqlDriverDB, createTableMssql, dropTableMssql,
		DoBenchDriverFind, t)
}

func BenchmarkMssqlNoCacheInsert(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	DoBenchInsert(engine, t)
}

func BenchmarkMssqlNoCacheFind(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	DoBenchFind(engine, t)
}

func BenchmarkMssqlNoCacheFindPtr(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	DoBenchFindPtr(engine, t)
}

func BenchmarkMssqlCacheInsert(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())

	DoBenchInsert(engine, t)
}

func BenchmarkMssqlCacheFind(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())

	DoBenchFind(engine, t)
}

func BenchmarkMssqlCacheFindPtr(t *testing.B) {
	engine, err := newMssqlEngine()
	if err != nil {
		t.Error(err)
		return
	}
	defer engine.Close()

	engine.SetDefaultCacher(NewCacher())

	DoBenchFindPtr(engine, t)
}
