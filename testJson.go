package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

type JsonProfile struct {
	Name string
	Age  int
}

type JsonField struct {
	Id      int64
	Name    map[string]string `xorm:"json"`
	Indexes []int             `xorm:"json"`
	Profile JsonProfile       `xorm:"json"`
}

func testJsonField(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&JsonField{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&JsonField{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	js := &JsonField{
		Name: map[string]string{
			"test":  "test",
			"test2": "test2",
		},
		Indexes: []int{1, 3, 5},
		Profile: JsonProfile{
			Name: "lll",
			Age:  12,
		},
	}
	_, err = engine.Insert(js)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var j JsonField
	has, err := engine.Id(js.Id).Get(&j)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println("j:", j)

	if !has {
		err = errors.New("not exist")
		t.Error(err)
		panic(err)
	}

	if j.Profile.Name != "lll" || j.Profile.Age != 12 {
		err = errors.New("json unmarshal error")
		t.Error(err)
		panic(err)
	}

	var jss = make([]JsonField, 0)
	err = engine.Find(&jss)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(jss) == 0 {
		err = errors.New("not exist")
		t.Error(err)
		panic(err)
	}

	fmt.Println("jss:", jss)
}
