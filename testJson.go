package tests

import (
    "fmt"
	"errors"
	"testing"

	"github.com/go-xorm/xorm"
)

type JsonField struct {
	Id   int64
	Name map[string]string `xorm:"json"`
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
			"test": "test",
			"test2": "test2",
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