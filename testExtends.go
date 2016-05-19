package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

type tempUser2 struct {
	TempUser tempUser   `xorm:"extends"`
	Departname string
}

type tempUser3 struct {
	Temp       *tempUser `xorm:"extends"`
	Departname string
}

type tempUser4 struct {
	TempUser2 tempUser2 `xorm:"extends"`
}

type UserAndDetail struct {
	Userinfo   `xorm:"extends"`
	Userdetail `xorm:"extends"`
}

func testExtends(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&tempUser2{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&tempUser2{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu := &tempUser2{tempUser{0, "extends"}, "dev depart"}
	_, err = engine.Insert(tu)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu2 := &tempUser2{}
	_, err = engine.Get(tu2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu3 := &tempUser2{tempUser{0, "extends update"}, ""}
	_, err = engine.Id(tu2.TempUser.Id).Update(tu3)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.DropTables(&tempUser4{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&tempUser4{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu8 := &tempUser4{tempUser2{tempUser{0, "extends"}, "dev depart"}}
	_, err = engine.Insert(tu8)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu9 := &tempUser4{}
	_, err = engine.Get(tu9)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if tu9.TempUser2.TempUser.Username != tu8.TempUser2.TempUser.Username || tu9.TempUser2.Departname != tu8.TempUser2.Departname {
		err = errors.New(fmt.Sprintln("not equal for", tu8, tu9))
		t.Error(err)
		panic(err)
	}

	tu10 := &tempUser4{tempUser2{tempUser{0, "extends update"}, ""}}
	_, err = engine.Id(tu9.TempUser2.TempUser.Id).Update(tu10)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.DropTables(&tempUser3{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&tempUser3{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu4 := &tempUser3{&tempUser{0, "extends"}, "dev depart"}
	_, err = engine.Insert(tu4)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu5 := &tempUser3{}
	_, err = engine.Get(tu5)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if tu5.Temp == nil {
		err = errors.New("error get data extends")
		t.Error(err)
		panic(err)
	}
	if tu5.Temp.Id != 1 || tu5.Temp.Username != "extends" ||
		tu5.Departname != "dev depart" {
		err = errors.New("error get data extends")
		t.Error(err)
		panic(err)
	}

	tu6 := &tempUser3{&tempUser{0, "extends update"}, ""}
	_, err = engine.Id(tu5.Temp.Id).Update(tu6)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	users := make([]tempUser3, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users) != 1 {
		err = errors.New("error get data not 1")
		t.Error(err)
		panic(err)
	}

	var info UserAndDetail
	qt := engine.Quote
	engine.Update(&Userinfo{Detail: Userdetail{Id: 1}})
	ui := engine.TableMapper.Obj2Table("Userinfo")
	ud := engine.TableMapper.Obj2Table("Userdetail")
	uiid := engine.TableMapper.Obj2Table("Id")
	udid := "detail_id"
	sql := fmt.Sprintf("select * from %s, %s where %s.%s = %s.%s",
		qt(ui), qt(ud), qt(ui), qt(udid), qt(ud), qt(uiid))
	b, err := engine.Sql(sql).Get(&info)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !b {
		err = errors.New("should has lest one record")
		t.Error(err)
		panic(err)
	}
	if info.Userinfo.Uid == 0 || info.Userdetail.Id == 0 {
		err = errors.New("all of the id should has value")
		t.Error(err)
		panic(err)
	}
	fmt.Println(info)

	fmt.Println("----join--info2")
	var info2 UserAndDetail
	b, err = engine.Table(&Userinfo{}).Join("LEFT", qt(ud), qt(ui)+"."+qt("detail_id")+" = "+qt(ud)+"."+qt(uiid)).Get(&info2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !b {
		err = errors.New("should has lest one record")
		t.Error(err)
		panic(err)
	}
	if info2.Userinfo.Uid == 0 || info2.Userdetail.Id == 0 {
		err = errors.New("all of the id should has value")
		t.Error(err)
		panic(err)
	}
	fmt.Println(info2)

	fmt.Println("----join--infos2")
	var infos2 = make([]UserAndDetail, 0)
	err = engine.Table(&Userinfo{}).Join("LEFT", qt(ud), qt(ui)+"."+qt("detail_id")+" = "+qt(ud)+"."+qt(uiid)).Find(&infos2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(infos2)

	testExtends2(engine, t)
	testExtends3(engine, t)
	testExtends4(engine, t)
}

type MessageBase struct {
	Id         int64     `xorm:"int(11) pk autoincr"`
	TypeId     int64     `xorm:"int(11) notnull"`
}

type Message struct {
	MessageBase `xorm:"extends"`
	Title      string    `xorm:"varchar(100) notnull"`
	Content    string    `xorm:"text notnull"`
	Uid        int64     `xorm:"int(11) notnull"`
	ToUid      int64     `xorm:"int(11) notnull"`
}

type MessageUser struct {
	Id   int64
	Name string
}

type MessageType struct {
	Id   int64
	Name string
}

type MessageExtend3 struct {
	Message  `xorm:"extends"`
	Sender   MessageUser `xorm:"extends"`
	Receiver MessageUser `xorm:"extends"`
	Type     MessageType `xorm:"extends"`
}

type MessageExtend4 struct {
	Message     `xorm:"extends"`
	MessageUser `xorm:"extends"`
	MessageType `xorm:"extends"`
}


func newExtendsTest(engine *xorm.Engine, t *testing.T) extendsTest {
	var ret extendsTest

	err := engine.DropTables(&Message{}, &MessageUser{}, &MessageType{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&Message{}, &MessageUser{}, &MessageType{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	ret.msgtype = MessageType{Name: "type"}
	ret.sender = MessageUser{Name: "sender"}
	ret.receiver = MessageUser{Name: "receiver"}

	_, err = engine.Insert(&ret.msgtype, &ret.sender, &ret.receiver)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	ret.message = Message{
		MessageBase: MessageBase{
			TypeId: ret.msgtype.Id,
		},
		Title:   "test",
		Content: "test",
		Uid:     ret.sender.Id,
		ToUid:   ret.receiver.Id,
	}

	_, err = engine.Insert(&ret.message)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	return ret
}

const (
	NeedSender = (1 << iota)
	NeedReceiver
	NeedType
	NeedAll = (NeedSender | NeedReceiver | NeedType)
)

func (e extendsTest) query(
	engine *xorm.Engine, t *testing.T, result interface{}, flags int) {
	mapper := engine.TableMapper.Obj2Table

	sess := engine.Table(mapper("Message"))

	if (flags & NeedSender) != 0 {
		sess.Join("LEFT", []string{mapper("MessageUser"), "sender"},
			fmt.Sprintf("`sender`.`%s` = `%s`.`%s`",
				mapper("Id"),
				mapper("Message"),
				mapper("Uid")))
	}

	if (flags & NeedReceiver) != 0 {
		sess.Join("LEFT", []string{mapper("MessageUser"), "receiver"},
			fmt.Sprintf("`receiver`.`%s` = `%s`.`%s`",
				mapper("Id"),
				mapper("Message"),
				mapper("ToUid")))
	}

	if (flags & NeedType) != 0 {
		sess.Join("LEFT", []string{mapper("MessageType"), "type"},
			fmt.Sprintf("`type`.`%s` = `%s`.`%s`",
				mapper("Id"),
				mapper("Message"),
				mapper("TypeId")))
	}

	if err := sess.Find(result); err != nil {
		t.Error(err)
		panic(err)
	}
}

func testExtends2(engine *xorm.Engine, t *testing.T) {
	e := newExtendsTest(engine, t)

	list := make([]Message, 0)

	e.query(engine, t, &list, NeedAll)

	if len(list) != 1 {
		err := errors.New(fmt.Sprintln("should have 1 message, got", len(list)))
		t.Error(err)
		panic(err)
	}

	if list[0] != e.message {
		err := errors.New(fmt.Sprintln("should message equal", list[0], e.message))
		t.Error(err)
		panic(err)
	}
}

func testExtends3(engine *xorm.Engine, t *testing.T) {
	e := newExtendsTest(engine, t)

	list := make([]MessageExtend3, 0)

	e.query(engine, t, &list, NeedAll)

	if len(list) != 1 {
		err := errors.New(fmt.Sprintln("should have 1 message, got", len(list)))
		t.Error(err)
		panic(err)
	}

	if list[0].Message != e.message {
		err := errors.New(fmt.Sprintln("should message equal", list[0].Message, e.message))
		t.Error(err)
		panic(err)
	}

	if list[0].Sender != e.sender {
		err := errors.New(fmt.Sprintln("should sender equal", list[0].Sender, e.sender))
		t.Error(err)
		panic(err)
	}

	if list[0].Receiver != e.receiver {
		err := errors.New(
			fmt.Sprintln("should receiver equal", list[0].Receiver, e.receiver))
		t.Error(err)
		panic(err)
	}

	if list[0].Type != e.msgtype {
		err := errors.New(fmt.Sprintln("should msgtype equal", list[0].Type, e.msgtype))
		t.Error(err)
		panic(err)
	}
}

func testExtends4(engine *xorm.Engine, t *testing.T) {
	e := newExtendsTest(engine, t)

	list := make([]MessageExtend4, 0)

	e.query(engine, t, &list, NeedSender | NeedType)

	if len(list) != 1 {
		err := errors.New(fmt.Sprintln("should have 1 message, got", len(list)))
		t.Error(err)
		panic(err)
	}

	if list[0].Message != e.message {
		err := errors.New(fmt.Sprintln("should message equal", list[0].Message, e.message))
		t.Error(err)
		panic(err)
	}

	if list[0].MessageUser != e.sender {
		err := errors.New(
			fmt.Sprintln("should sender equal", list[0].MessageUser, e.sender))
		t.Error(err)
		panic(err)
	}

	if list[0].MessageType != e.msgtype {
		err := errors.New(
			fmt.Sprintln("should msgtype equal", list[0].MessageType, e.msgtype))
		t.Error(err)
		panic(err)
	}
}
