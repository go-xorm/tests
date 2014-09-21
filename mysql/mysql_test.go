package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	. ".."
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

/*
CREATE DATABASE IF NOT EXISTS xorm_test CHARACTER SET
utf8 COLLATE utf8_general_ci;
*/

func getDSN(dbName string) string {
	user := os.Getenv("GOSQLTEST_MYSQL_USER")
	if user == "" {
		user = "root"
	}
	pass, _ := GetEnvOk("GOSQLTEST_MYSQL_PASS")
	return fmt.Sprintf("%s:%s@/%s?charset=utf8", user, pass, dbName)
}

func TestMysql(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", getDSN("xorm_test"))
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.ShowSQL = ShowTestSql
	engine.ShowErr = ShowTestSql
	engine.ShowWarn = ShowTestSql
	engine.ShowDebug = ShowTestSql

	BaseTestAll(engine, t)
	BaseTestAllSnakeMapper(engine, t)
	BaseTestAll2(engine, t)
	BaseTestAll3(engine, t)
}

func TestMysqlSameMapper(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", getDSN("xorm_test1"))
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.ShowSQL = ShowTestSql
	engine.ShowErr = ShowTestSql
	engine.ShowWarn = ShowTestSql
	engine.ShowDebug = ShowTestSql
	engine.SetMapper(core.SameMapper{})

	BaseTestAll(engine, t)
	BaseTestAllSameMapper(engine, t)
	BaseTestAll2(engine, t)
	BaseTestAll3(engine, t)
}

func TestMysqlWithCache(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", getDSN("xorm_test2"))
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())
	engine.ShowSQL = ShowTestSql
	engine.ShowErr = ShowTestSql
	engine.ShowWarn = ShowTestSql
	engine.ShowDebug = ShowTestSql

	BaseTestAll(engine, t)
	BaseTestAllSnakeMapper(engine, t)
	BaseTestAll2(engine, t)
}

func TestMysqlWithCacheSameMapper(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", getDSN("xorm_test3"))
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetMapper(core.SameMapper{})
	engine.SetDefaultCacher(NewCacher())
	engine.ShowSQL = ShowTestSql
	engine.ShowErr = ShowTestSql
	engine.ShowWarn = ShowTestSql
	engine.ShowDebug = ShowTestSql

	BaseTestAll(engine, t)
	BaseTestAllSameMapper(engine, t)
	BaseTestAll2(engine, t)
}

func newMysqlEngine() (*xorm.Engine, error) {
	return xorm.NewEngine("mysql", getDSN("xorm_test"))
}

func mysqlDdlImport() error {
	engine, err := xorm.NewEngine("mysql", getDSN(""))
	if err != nil {
		return err
	}
	engine.ShowSQL = ShowTestSql
	engine.ShowErr = ShowTestSql
	engine.ShowWarn = ShowTestSql
	engine.ShowDebug = ShowTestSql

	sqlResults, _ := engine.ImportFile("../testdata/mysql_ddl.sql")
	engine.LogDebug("sql results: %v", sqlResults)
	engine.Close()
	return nil
}

func newMysqlDriverDB() (*sql.DB, error) {
	return sql.Open("mysql", getDSN("xorm_test"))
}

func BenchmarkMysqlDriverInsert(t *testing.B) {
	DoBenchDriver(newMysqlDriverDB, CreateTableMySql, DropTableMySql,
		DoBenchDriverInsert, t)
}

func BenchmarkMysqlDriverFind(t *testing.B) {
	DoBenchDriver(newMysqlDriverDB, CreateTableMySql, DropTableMySql,
		DoBenchDriverFind, t)
}

func BenchmarkMysqlNoCacheInsert(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchInsert(engine, t)
}

func BenchmarkMysqlNoCacheFind(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchFind(engine, t)
}

func BenchmarkMysqlNoCacheFindPtr(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	//engine.ShowSQL = true
	DoBenchFindPtr(engine, t)
}

func BenchmarkMysqlCacheInsert(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())

	DoBenchInsert(engine, t)
}

func BenchmarkMysqlCacheFind(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())

	DoBenchFind(engine, t)
}

func BenchmarkMysqlCacheFindPtr(t *testing.B) {
	engine, err := newMysqlEngine()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewCacher())

	DoBenchFindPtr(engine, t)
}
