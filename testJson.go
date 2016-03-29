package tests

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-xorm/xorm"
)

type JsonProfile struct {
	Name string
	Age  int
}

type JsonField struct {
	Id       int64
	Name     map[string]string `xorm:"json"`
	Indexes  []int             `xorm:"json"`
	Profile  JsonProfile       `xorm:"json"`
	Profile2 *JsonProfile      `xorm:"json"`
	Name2    map[string]string
	Indexes2 []int
	//Profile3 JsonProfile
	//Profile4 *JsonProfile
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
		Profile2: &JsonProfile{
			Name: "lll",
			Age:  12,
		},
		Name2: map[string]string{
			"test":  "test",
			"test2": "test2",
		},
		Indexes2: []int{1, 3, 5},
		/*Profile3: JsonProfile{
			Name: "lll",
			Age:  12,
		},
		Profile4: &JsonProfile{
			Name: "lll",
			Age:  12,
		},*/
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

	if !reflect.DeepEqual(js, &j) {
		err = fmt.Errorf("%v is not equal %v", *js, j)
		t.Error(err)
		panic(err)
	}

	var jss = make([]JsonField, 0)
	err = engine.Find(&jss)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(jss) != 1 {
		err = errors.New("not exist")
		t.Error(err)
		panic(err)
	}

	if !reflect.DeepEqual(js, &jss[0]) {
		err = fmt.Errorf("%v is not equal %v", *js, j)
		t.Error(err)
		panic(err)
	}
}
