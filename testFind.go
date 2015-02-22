package tests

import (
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

func find(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}

	users2 := make([]Userinfo, 0)
	userinfo := engine.TableMapper.Obj2Table("Userinfo")
	err = engine.Sql("select * from " + engine.Quote(userinfo)).Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func find2(engine *xorm.Engine, t *testing.T) {
	users := make([]*Userinfo, 0)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}

func findMap(engine *xorm.Engine, t *testing.T) {
	users := make(map[int64]Userinfo)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}

func findMap2(engine *xorm.Engine, t *testing.T) {
	users := make(map[int64]*Userinfo)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for id, user := range users {
		fmt.Println(id, user)
	}
}
