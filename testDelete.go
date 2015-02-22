package tests

import (
	"testing"

	"github.com/go-xorm/xorm"
)

func testDelete(engine *xorm.Engine, t *testing.T) {
	user := Userinfo{Uid: 1}
	cnt, err := engine.Delete(&user)
	if err != nil {
		t.Fatal("delete failed:", err)
	}
	if cnt != 1 {
		t.Fatal("delete failed: deleted 0 rows")
	}

	user.Uid = 0
	user.IsMan = true
	has, err := engine.Id(3).Get(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if has {
		//var tt time.Time
		//user.Created = tt
		cnt, err := engine.Id(3).Delete(new(Userinfo))
		if err != nil {
			t.Fatal("delete failed:", err)
		}
		if cnt != 1 {
			t.Fatal("delete failed: deleted 0 rows")
		}
	}
}
