package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
)

func testExtends(engine *xorm.Engine, t *testing.T) {
	testExtends1(engine, t)
	testExtends2(engine, t)
	testExtends3(engine, t)
	testExtends4(engine, t)
	testExtends5(engine, t)
	testExtends6(engine, t)
	testExtends7(engine, t)
	testExtends8(engine, t)
	testExtends9(engine, t)
}

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

func testExtends1(engine *xorm.Engine, t *testing.T) {
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

type MessageBaseAndTitle struct {
	MessageBase `xorm:"extends"`
	Title string `xorm:"varchar(100) notnull"`
}

type MessageUserName struct {
	Name string
}

type MessageExtend3 struct {
	Message  `xorm:"extends"`
	Sender   MessageUser `xorm:"extends"`
	Receiver MessageUser `xorm:"extends"`
	Type     MessageType `xorm:"extends"`
}

type MessageExtend4 struct {
	Message  `xorm:"extends"`
	Sender   MessageUserName `xorm:"extends"`
	Receiver MessageUserName `xorm:"extends"`
}

type MessageExtend5 struct {
	Message     `xorm:"extends"`
	MessageUser `xorm:"extends"`
	MessageType `xorm:"extends"`
}

type MessageExtend6 struct {
	Message `xorm:"extends"`
	MyUser MessageUser `xorm:"extends"`
	MyType MessageType `xorm:"extends"`
}

type MessageExtend7 struct {
	Message     `xorm:"extends"`
	MyUser1 MessageUser `xorm:"extends"`
	MyUser2 MessageUser `xorm:"extends"`
}

type MessageExtend8 struct {
	Message     MessageBaseAndTitle `xorm:"extends"`
	MessageUser MessageUser         `xorm:"extends"`
}

type MessageExtend9 struct {
	MessageBaseAndTitle `xorm:"extends"`
	MessageUser         `xorm:"extends"`
	Content string
}

type extendsTest struct {
	message  Message
	msgtype  MessageType
	sender   MessageUser
	receiver MessageUser
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

	ret.msgtype = MessageType{Name: "type_name"}
	ret.sender = MessageUser{Name: "sender_name"}
	ret.receiver = MessageUser{Name: "receiver_name"}

	_, err = engine.Insert(&ret.msgtype, &ret.sender, &ret.receiver)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	ret.message = Message{
		MessageBase: MessageBase{
			TypeId: ret.msgtype.Id,
		},
		Title:   "test_title",
		Content: "test_content",
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
	AllowAmbiguous = (1 << iota)
	NeedSender
	NeedReceiver
	NeedType
	NeedAll = (NeedSender | NeedReceiver | NeedType)
)

func (e extendsTest) query(
	engine *xorm.Engine, t *testing.T, result interface{}, flags int) {
	mapper := engine.TableMapper.Obj2Table

	sess := engine.Table(mapper("Message"))

	if (flags & AllowAmbiguous) == 0 {
		sess.Statement.AllowAmbiguous = false
	}

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

	if list[0].Sender.Name != e.sender.Name {
		err := errors.New(fmt.Sprintln("should sender equal", list[0].Sender, e.sender))
		t.Error(err)
		panic(err)
	}

	if list[0].Receiver.Name != e.receiver.Name {
		err := errors.New(
			fmt.Sprintln("should receiver equal", list[0].Receiver, e.receiver))
		t.Error(err)
		panic(err)
	}
}

func testExtends5(engine *xorm.Engine, t *testing.T) {
	e := newExtendsTest(engine, t)

	list := make([]MessageExtend5, 0)

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

func testExtends6(engine *xorm.Engine, t *testing.T) {
	e := newExtendsTest(engine, t)

	list := make([]MessageExtend6, 0)

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

	if list[0].MyUser != e.sender {
		err := errors.New(fmt.Sprintln("should sender equal", list[0].MyUser, e.sender))
		t.Error(err)
		panic(err)
	}

	if list[0].MyType != e.msgtype {
		err := errors.New(fmt.Sprintln("should msgtype equal", list[0].MyType, e.msgtype))
		t.Error(err)
		panic(err)
	}
}

func testExtends7(engine *xorm.Engine, t *testing.T) {
	e := newExtendsTest(engine, t)

	list := make([]MessageExtend7, 0)

	e.query(engine, t, &list, NeedAll | AllowAmbiguous)

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

	if !((list[0].MyUser1 == e.sender && list[0].MyUser2 == e.receiver) ||
		(list[0].MyUser2 == e.sender && list[0].MyUser1 == e.receiver)) {
		err := errors.New(fmt.Sprintln("should sender receiver equal",
			list[0].MyUser1, list[0].MyUser2, e.sender, e.receiver))
		t.Error(err)
		panic(err)
	}
}

func testExtends8(engine *xorm.Engine, t *testing.T) {
	e := newExtendsTest(engine, t)

	list := make([]MessageExtend8, 0)

	e.query(engine, t, &list, NeedSender)

	if len(list) != 1 {
		err := errors.New(fmt.Sprintln("should have 1 message, got", len(list)))
		t.Error(err)
		panic(err)
	}

	if list[0].Message.MessageBase != e.message.MessageBase ||
		list[0].Message.Title != e.message.Title {
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
}

func testExtends9(engine *xorm.Engine, t *testing.T) {
	e := newExtendsTest(engine, t)

	list := make([]MessageExtend9, 0)

	e.query(engine, t, &list, NeedSender | AllowAmbiguous)

	if len(list) != 1 {
		err := errors.New(fmt.Sprintln("should have 1 message, got", len(list)))
		t.Error(err)
		panic(err)
	}

	if list[0].MessageBaseAndTitle.MessageBase != e.message.MessageBase ||
		list[0].MessageBaseAndTitle.Title != e.message.Title {
		err := errors.New(fmt.Sprintln("should message equal",
			list[0].MessageBaseAndTitle, e.message))
		t.Error(err)
		panic(err)
	}

	if list[0].Content != e.message.Content {
		err := errors.New(fmt.Sprintln("should message content equal",
			list[0].Content, e.message.Content))
		t.Error(err)
		panic(err)
	}

	if list[0].MessageUser != e.sender {
		err := errors.New(
			fmt.Sprintln("should sender equal", list[0].MessageUser, e.sender))
		t.Error(err)
		panic(err)
	}
}
