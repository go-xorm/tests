package tests

import (
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

func count(engine *xorm.Engine, t *testing.T) {
	user := Userinfo{Departname: "dev"}
	total, err := engine.Count(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Printf("Total %d records!!!\n", total)
}
