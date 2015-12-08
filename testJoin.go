package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/xorm"
)

func joinSameMapper(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Join("LEFT", "`Userdetail`", "`Userinfo`.`(id)`=`Userdetail`.`Id`").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func join(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Join("LEFT", "userdetail", "userinfo.id=userdetail.id").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	join2(engine, t)
	join3(engine, t)
	joinCount(engine, t)
	joinCount2(engine, t)
}

func join2(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Join("LEFT", "userdetail", "userinfo.id=userdetail.id").Find(&users,
		&Userinfo{Uid: 1})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func join3(engine *xorm.Engine, t *testing.T) {
	_, err := engine.Join("LEFT", "userdetail", "userinfo.id=userdetail.id").Get(&Userinfo{Uid: 1})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func joinCount(engine *xorm.Engine, t *testing.T) {
	count, err := engine.Join("LEFT", "userdetail", "userinfo.id=userdetail.id").Count(&Userinfo{Uid: 1})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(count)
}

type History struct {
	Rid       int64
	Uid       int64
	DeletedAt time.Time `xorm:"deleted"`
}

type Resource struct {
	Rid int64
}

type NewUser struct {
	Uid int64
}

func (NewUser) TableName() string {
	return "user"
}

func joinCount2(engine *xorm.Engine, t *testing.T) {
	err := engine.Sync2(new(History), new(Resource), new(NewUser))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var where = "history.deleted_at > '0000-00-00 00:00:00'"
	count, err := engine.Table("history").Join("LEFT", "resource", "resource.rid=history.rid").Join("LEFT", "user", "user.uid=history.uid").Where(where).Count(new(History))
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(count)
}
