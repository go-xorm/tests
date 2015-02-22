package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

type Condi map[string]interface{}

type UpdateAllCols struct {
	Id     int64
	Bool   bool
	String string
	Ptr    *string
}

type UpdateMustCols struct {
	Id     int64
	Bool   bool
	String string
}

type UpdateIncr struct {
	Id  int64
	Cnt int
}

func update(engine *xorm.Engine, t *testing.T) {
	// update by id
	user := Userinfo{Username: "xxx", Height: 1.2}
	cnt, err := engine.Id(4).Update(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	condi := Condi{"username": "zzz", "departname": ""}
	cnt, err = engine.Table(&user).Id(4).Update(&condi)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	cnt, err = engine.Update(&Userinfo{Username: "yyy"}, &user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	total, err := engine.Count(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != total {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	err = engine.StoreEngine("Innodb").Sync2(&Article{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	defer func() {
		err = engine.DropTables(&Article{})
		if err != nil {
			t.Error(err)
			panic(err)
		}
	}()

	a := &Article{0, "1", "2", "3", "4", "5", 2}
	cnt, err = engine.Insert(a)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New(fmt.Sprintf("insert not returned 1 but %d", cnt))
		t.Error(err)
		panic(err)
	}

	if a.Id == 0 {
		err = errors.New("insert returned id is 0")
		t.Error(err)
		panic(err)
	}

	cnt, err = engine.Id(a.Id).Update(&Article{Name: "6"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New(fmt.Sprintf("insert not returned 1 but %d", cnt))
		t.Error(err)
		panic(err)
		return
	}

	var s = "test"

	col1 := &UpdateAllCols{Ptr: &s}
	err = engine.Sync(col1)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Insert(col1)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	col2 := &UpdateAllCols{col1.Id, true, "", nil}
	_, err = engine.Id(col2.Id).AllCols().Update(col2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	col3 := &UpdateAllCols{}
	has, err := engine.Id(col2.Id).Get(col3)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		err = errors.New(fmt.Sprintf("cannot get id %d", col2.Id))
		t.Error(err)
		panic(err)
		return
	}

	if *col2 != *col3 {
		err = errors.New(fmt.Sprintf("col2 should eq col3"))
		t.Error(err)
		panic(err)
		return
	}

	{

		col1 := &UpdateMustCols{}
		err = engine.Sync(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		_, err = engine.Insert(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		col2 := &UpdateMustCols{col1.Id, true, ""}
		boolStr := engine.ColumnMapper.Obj2Table("Bool")
		stringStr := engine.ColumnMapper.Obj2Table("String")
		_, err = engine.Id(col2.Id).MustCols(boolStr, stringStr).Update(col2)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		col3 := &UpdateMustCols{}
		has, err := engine.Id(col2.Id).Get(col3)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		if !has {
			err = errors.New(fmt.Sprintf("cannot get id %d", col2.Id))
			t.Error(err)
			panic(err)
			return
		}

		if *col2 != *col3 {
			err = errors.New(fmt.Sprintf("col2 should eq col3"))
			t.Error(err)
			panic(err)
			return
		}
	}

	{

		col1 := &UpdateIncr{}
		err = engine.Sync(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		_, err = engine.Insert(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		cnt, err := engine.Id(col1.Id).Incr("cnt").Update(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if cnt != 1 {
			err = errors.New("update incr failed")
			t.Error(err)
			panic(err)
		}

		newCol := new(UpdateIncr)
		has, err := engine.Id(col1.Id).Get(newCol)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if !has {
			err = errors.New("has incr failed")
			t.Error(err)
			panic(err)
		}
		if 1 != newCol.Cnt {
			err = fmt.Errorf("incr failed %v %v %v", newCol.Cnt, newCol, col1)
			t.Error(err)
			panic(err)
		}
	}
}

func updateSameMapper(engine *xorm.Engine, t *testing.T) {
	// update by id
	user := Userinfo{Username: "xxx", Height: 1.2}
	cnt, err := engine.Id(4).Update(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	condi := Condi{"Username": "zzz", "Departname": ""}
	cnt, err = engine.Table(&user).Id(4).Update(&condi)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	cnt, err = engine.Update(&Userinfo{Username: "yyy"}, &user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	total, err := engine.Count(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != total {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	err = engine.Sync(&Article{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	defer func() {
		err = engine.DropTables(&Article{})
		if err != nil {
			t.Error(err)
			panic(err)
		}
	}()

	a := &Article{0, "1", "2", "3", "4", "5", 2}
	cnt, err = engine.Insert(a)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New(fmt.Sprintf("insert not returned 1 but %d", cnt))
		t.Error(err)
		panic(err)
	}

	if a.Id == 0 {
		err = errors.New("insert returned id is 0")
		t.Error(err)
		panic(err)
	}

	cnt, err = engine.Id(a.Id).Update(&Article{Name: "6"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New(fmt.Sprintf("insert not returned 1 but %d", cnt))
		t.Error(err)
		panic(err)
		return
	}

	col1 := &UpdateAllCols{}
	err = engine.Sync(col1)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Insert(col1)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	col2 := &UpdateAllCols{col1.Id, true, "", nil}
	_, err = engine.Id(col2.Id).AllCols().Update(col2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	col3 := &UpdateAllCols{}
	has, err := engine.Id(col2.Id).Get(col3)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		err = errors.New(fmt.Sprintf("cannot get id %d", col2.Id))
		t.Error(err)
		panic(err)
		return
	}

	if *col2 != *col3 {
		err = errors.New(fmt.Sprintf("col2 should eq col3"))
		t.Error(err)
		panic(err)
		return
	}

	{
		col1 := &UpdateMustCols{}
		err = engine.Sync(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		_, err = engine.Insert(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		col2 := &UpdateMustCols{col1.Id, true, ""}
		boolStr := engine.ColumnMapper.Obj2Table("Bool")
		stringStr := engine.ColumnMapper.Obj2Table("String")
		_, err = engine.Id(col2.Id).MustCols(boolStr, stringStr).Update(col2)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		col3 := &UpdateMustCols{}
		has, err := engine.Id(col2.Id).Get(col3)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		if !has {
			err = errors.New(fmt.Sprintf("cannot get id %d", col2.Id))
			t.Error(err)
			panic(err)
			return
		}

		if *col2 != *col3 {
			err = errors.New(fmt.Sprintf("col2 should eq col3"))
			t.Error(err)
			panic(err)
			return
		}
	}

	{

		col1 := &UpdateIncr{}
		err = engine.Sync(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		_, err = engine.Insert(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}

		cnt, err := engine.Id(col1.Id).Incr("`Cnt`").Update(col1)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if cnt != 1 {
			err = errors.New("update incr failed")
			t.Error(err)
			panic(err)
		}

		newCol := new(UpdateIncr)
		has, err := engine.Id(col1.Id).Get(newCol)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if !has {
			err = errors.New("has incr failed")
			t.Error(err)
			panic(err)
		}
		if 1 != newCol.Cnt {
			err = errors.New("incr failed")
			t.Error(err)
			panic(err)
		}
	}
}
