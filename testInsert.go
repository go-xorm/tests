package tests

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/xorm"
)

func insert(engine *xorm.Engine, t *testing.T) {
	user := Userinfo{0, "xiaolunwen", "dev", "lunny", time.Now(),
		Userdetail{Id: 1}, 1.78, []byte{1, 2, 3}, true}
	cnt, err := engine.Insert(&user)
	fmt.Println(user.Uid)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
	if user.Uid <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	user.Uid = 0
	cnt, err = engine.Insert(&user)
	if err == nil {
		err = errors.New("insert failed but no return error")
		t.Error(err)
		panic(err)
	}
	if cnt != 0 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	testInsertCreated(engine, t)
}

func insertAutoIncr(engine *xorm.Engine, t *testing.T) {
	// auto increment insert
	user := Userinfo{Username: "xiaolunwen2", Departname: "dev", Alias: "lunny", Created: time.Now(),
		Detail: Userdetail{Id: 1}, Height: 1.78, Avatar: []byte{1, 2, 3}, IsMan: true}
	cnt, err := engine.Insert(&user)
	fmt.Println(user.Uid)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
	if user.Uid <= 0 {
		t.Error(errors.New("not return id error"))
	}
}

type DefaultInsert struct {
	Id      int64
	Status  int `xorm:"default -1"`
	Name    string
	Updated time.Time `xorm:"updated"`
}

func testInsertDefault(engine *xorm.Engine, t *testing.T) {
	di := new(DefaultInsert)
	err := engine.Sync2(di)
	if err != nil {
		t.Error(err)
	}

	_, err = engine.Omit(engine.ColumnMapper.Obj2Table("Status")).Insert(&DefaultInsert{Name: "test"})
	if err != nil {
		t.Error(err)
	}

	has, err := engine.Desc("(id)").Get(di)
	if err != nil {
		t.Error(err)
	}
	if !has {
		err = errors.New("error with no data")
		t.Error(err)
		panic(err)
	}
	if di.Status != -1 {
		err = errors.New("inserted error data")
		t.Error(err)
		panic(err)
	}
}

type CreatedInsert struct {
	Id      int64
	Created time.Time `xorm:"created"`
}

type CreatedInsert2 struct {
	Id      int64
	Created int64 `xorm:"created"`
}

type CreatedInsert3 struct {
	Id      int64
	Created int `xorm:"created bigint"`
}

type CreatedInsert4 struct {
	Id      int64
	Created int `xorm:"created"`
}

type CreatedInsert5 struct {
	Id      int64
	Created time.Time `xorm:"created bigint"`
}

func testInsertCreated(engine *xorm.Engine, t *testing.T) {
	di := new(CreatedInsert)
	err := engine.Sync2(di)
	if err != nil {
		t.Fatal(err)
	}
	ci := &CreatedInsert{}
	_, err = engine.Insert(ci)
	if err != nil {
		t.Fatal(err)
	}

	has, err := engine.Desc("(id)").Get(di)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci.Created.Unix() != di.Created.Unix() {
		t.Fatal("should equal:", ci, di)
	}
	fmt.Println("ci:", ci, "di:", di)

	di2 := new(CreatedInsert2)
	err = engine.Sync2(di2)
	if err != nil {
		t.Fatal(err)
	}
	ci2 := &CreatedInsert2{}
	_, err = engine.Insert(ci2)
	if err != nil {
		t.Fatal(err)
	}
	has, err = engine.Desc("(id)").Get(di2)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci2.Created != di2.Created {
		t.Fatal("should equal:", ci2, di2)
	}
	fmt.Println("ci2:", ci2, "di2:", di2)

	di3 := new(CreatedInsert3)
	err = engine.Sync2(di3)
	if err != nil {
		t.Fatal(err)
	}
	ci3 := &CreatedInsert3{}
	_, err = engine.Insert(ci3)
	if err != nil {
		t.Fatal(err)
	}
	has, err = engine.Desc("(id)").Get(di3)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci3.Created != di3.Created {
		t.Fatal("should equal:", ci3, di3)
	}
	fmt.Println("ci3:", ci3, "di3:", di3)

	di4 := new(CreatedInsert4)
	err = engine.Sync2(di4)
	if err != nil {
		t.Fatal(err)
	}
	ci4 := &CreatedInsert4{}
	_, err = engine.Insert(ci4)
	if err != nil {
		t.Fatal(err)
	}
	has, err = engine.Desc("(id)").Get(di4)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci4.Created != di4.Created {
		t.Fatal("should equal:", ci4, di4)
	}
	fmt.Println("ci4:", ci4, "di4:", di4)

	di5 := new(CreatedInsert5)
	err = engine.Sync2(di5)
	if err != nil {
		t.Fatal(err)
	}
	ci5 := &CreatedInsert5{}
	_, err = engine.Insert(ci5)
	if err != nil {
		t.Fatal(err)
	}
	has, err = engine.Desc("(id)").Get(di5)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Fatal(xorm.ErrNotExist)
	}
	if ci5.Created.Unix() != di5.Created.Unix() {
		t.Fatal("should equal:", ci5, di5)
	}
	fmt.Println("ci5:", ci5, "di5:", di5)
}

func insertMulti(engine *xorm.Engine, t *testing.T) {
	users := []Userinfo{
		{Username: "xlw", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		{Username: "xlw2", Departname: "dev", Alias: "lunny3", Created: time.Now()},
		{Username: "xlw11", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		{Username: "xlw22", Departname: "dev", Alias: "lunny3", Created: time.Now()},
	}
	cnt, err := engine.Insert(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != int64(len(users)) {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	users2 := []*Userinfo{
		&Userinfo{Username: "1xlw", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		&Userinfo{Username: "1xlw2", Departname: "dev", Alias: "lunny3", Created: time.Now()},
		&Userinfo{Username: "1xlw11", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		&Userinfo{Username: "1xlw22", Departname: "dev", Alias: "lunny3", Created: time.Now()},
	}

	cnt, err = engine.Insert(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != int64(len(users2)) {
		err = errors.New(fmt.Sprintf("insert not returned %v", len(users2)))
		t.Error(err)
		panic(err)
		return
	}
}

func insertTwoTable(engine *xorm.Engine, t *testing.T) {
	userdetail := Userdetail{ /*Id: 1, */ Intro: "I'm a very beautiful women.", Profile: "sfsaf"}
	userinfo := Userinfo{Username: "xlw3", Departname: "dev", Alias: "lunny4", Created: time.Now(), Detail: userdetail}

	cnt, err := engine.Insert(&userinfo, &userdetail)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if userinfo.Uid <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	if userdetail.Id <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	if cnt != 2 {
		err = errors.New("insert not returned 2")
		t.Error(err)
		panic(err)
		return
	}
}
