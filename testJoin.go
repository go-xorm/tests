package tests

import (
	"testing"

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
}
