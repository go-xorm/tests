package tests

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

const (
	CreateTableMySql = "CREATE TABLE IF NOT EXISTS `big_struct` (`id` BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL, `name` VARCHAR(255) NULL, `title` VARCHAR(255) NULL, `age` VARCHAR(255) NULL, `alias` VARCHAR(255) NULL, `nick_name` VARCHAR(255) NULL);"
	DropTableMySql   = "DROP TABLE IF EXISTS `big_struct`;"
)

var ShowTestSql bool = true

/*
CREATE TABLE `userinfo` (
    `id` INT(10) NULL AUTO_INCREMENT,
    `username` VARCHAR(64) NULL,
    `departname` VARCHAR(64) NULL,
    `created` DATE NULL,
    PRIMARY KEY (`uid`)
);
CREATE TABLE `userdeatail` (
    `id` INT(10) NULL,
    `intro` TEXT NULL,
    `profile` TEXT NULL,
    PRIMARY KEY (`uid`)
);
*/

type Userinfo struct {
	Uid        int64  `xorm:"id pk not null autoincr"`
	Username   string `xorm:"unique"`
	Departname string
	Alias      string `xorm:"-"`
	Created    time.Time
	Detail     Userdetail `xorm:"detail_id int(11)"`
	Height     float64
	Avatar     []byte
	IsMan      bool
}

type Userdetail struct {
	Id      int64
	Intro   string `xorm:"text"`
	Profile string `xorm:"varchar(2000)"`
}

type Picture struct {
	Id          int64
	Url         string `xorm:"unique"` //image's url
	Title       string
	Description string
	Created     time.Time `xorm:"created"`
	ILike       int
	PageView    int
	From_url    string
	Pre_url     string `xorm:"unique"` //pre view image's url
	Uid         int64
}

type Numeric struct {
	Numeric float64 `xorm:"numeric(26,2)"`
}

func NewCacher() core.Cacher {
	return xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Hour, 10000)
}

func testQuery(engine *xorm.Engine, t *testing.T) {
	sql := "select * from userinfo"
	results, err := engine.Query(sql)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(results)
}

func exec(engine *xorm.Engine, t *testing.T) {
	sql := "update userinfo set username=? where id=?"
	res, err := engine.Exec(sql, "xiaolun", 1)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(res)
}

func testQuerySameMapper(engine *xorm.Engine, t *testing.T) {
	sql := "select * from `Userinfo`"
	results, err := engine.Query(sql)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(results)
}

func execSameMapper(engine *xorm.Engine, t *testing.T) {
	sql := "update `Userinfo` set `Username`=? where (id)=?"
	res, err := engine.Exec(sql, "xiaolun", 1)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(res)
}

type Article struct {
	Id      int32  `xorm:"pk INT autoincr"`
	Name    string `xorm:"VARCHAR(45)"`
	Img     string `xorm:"VARCHAR(100)"`
	Aside   string `xorm:"VARCHAR(200)"`
	Desc    string `xorm:"VARCHAR(200)"`
	Content string `xorm:"TEXT"`
	Status  int8   `xorm:"TINYINT(4)"`
}

func where(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Where("(id) > ?", 2).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	err = engine.Where("(id) > ?", 2).And("(id) < ?", 10).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
}

func in(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.In("(id)", 7, 8, 9).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
	if len(users) != 3 {
		err = errors.New("in uses should be 7,8,9 total 3")
		t.Error(err)
		panic(err)
	}

	users = make([]Userinfo, 0)
	err = engine.In("(id)", []int{7, 8, 9}).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
	if len(users) != 3 {
		err = errors.New("in uses should be 7,8,9 total 3")
		t.Error(err)
		panic(err)
	}

	for _, user := range users {
		if user.Uid != 7 && user.Uid != 8 && user.Uid != 9 {
			err = errors.New("in uses should be 7,8,9 total 3")
			t.Error(err)
			panic(err)
		}
	}

	users = make([]Userinfo, 0)
	ids := []interface{}{7, 8, 9}
	department := engine.ColumnMapper.Obj2Table("Departname")
	err = engine.Where("`"+department+"` = ?", "dev").In("(id)", ids...).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	if len(users) != 3 {
		err = errors.New("in uses should be 7,8,9 total 3")
		t.Error(err)
		panic(err)
	}

	for _, user := range users {
		if user.Uid != 7 && user.Uid != 8 && user.Uid != 9 {
			err = errors.New("in uses should be 7,8,9 total 3")
			t.Error(err)
			panic(err)
		}
	}

	dev := engine.ColumnMapper.Obj2Table("Dev")

	err = engine.In("(id)", 1).In("(id)", 2).In(department, dev).Find(&users)

	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	cnt, err := engine.In("(id)", 4).Update(&Userinfo{Departname: "dev-"})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update records not 1")
		t.Error(err)
		panic(err)
	}

	user := new(Userinfo)
	has, err := engine.Id(4).Get(user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("get record not 1")
		t.Error(err)
		panic(err)
	}
	if user.Departname != "dev-" {
		err = errors.New("update not success")
		t.Error(err)
		panic(err)
	}

	cnt, err = engine.In("(id)", 4).Update(&Userinfo{Departname: "dev"})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update records not 1")
		t.Error(err)
		panic(err)
	}

	cnt, err = engine.In("(id)", 5).Delete(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("deleted records not 1")
		t.Error(err)
		panic(err)
	}
}

type Limit struct {
	Id      int64
	Name    string
	Updated time.Time `xorm:"updated"`
}

func limit(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Limit(2, 1).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	err = engine.DropTables(new(Limit))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Sync2(new(Limit))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Insert(&Limit{Name: "1"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	// TODO: support limit for update
	/*_, err = engine.Limit(1).Update(&Limit{Name: "2"})
	if err != nil {
		t.Error(err)
		panic(err)
	}*/
}

func order(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.OrderBy("id desc").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	users2 := make([]Userinfo, 0)
	err = engine.Asc("id", "username").Desc("height").Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users2)
}

func having(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.GroupBy("username").Having("username='xlw'").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	/*users = make([]Userinfo, 0)
	err = engine.Cols("id, username").GroupBy("username").Having("username='xlw'").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)*/
}

func orderSameMapper(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.OrderBy("(id) desc").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	users2 := make([]Userinfo, 0)
	err = engine.Asc("(id)", "Username").Desc("Height").Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users2)
}

func havingSameMapper(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.GroupBy("Username").Having(`"Username"='xlw'`).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
}

func tableOp(engine *xorm.Engine, t *testing.T) {
	user := Userinfo{Username: "tablexiao", Departname: "dev", Alias: "lunny", Created: time.Now()}
	tableName := fmt.Sprintf("user_%v", len(user.Username))
	cnt, err := engine.Table(tableName).Insert(&user)
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

	has, err := engine.Table(tableName).Get(&Userinfo{Username: "tablexiao"})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("Get has return false")
		t.Error(err)
		panic(err)
		return
	}

	users := make([]Userinfo, 0)
	err = engine.Table(tableName).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	id := user.Uid
	cnt, err = engine.Table(tableName).Id(id).Update(&Userinfo{Username: "tableda"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Table(tableName).Id(id).Delete(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testCharst(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables("user_charset")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Charset("utf8").Table("user_charset").CreateTable(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testStoreEngine(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables("user_store_engine")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.StoreEngine("InnoDB").Table("user_store_engine").CreateTable(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

type tempUser struct {
	Id       int64
	Username string
}

func testCols(engine *xorm.Engine, t *testing.T) {
	users := []Userinfo{}
	err := engine.Cols("id, username").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println(users)

	tmpUsers := []tempUser{}
	err = engine.NoCache().Table("userinfo").Cols("id, username").Find(&tmpUsers)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(tmpUsers)

	user := &Userinfo{Uid: 1, Alias: "", Height: 0}
	affected, err := engine.Cols("departname, height").Id(1).Update(user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println("===================", user, affected)
}

func testColsSameMapper(engine *xorm.Engine, t *testing.T) {
	users := []Userinfo{}
	err := engine.Cols("id, Username").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println(users)

	tmpUsers := []tempUser{}
	// TODO: should use cache
	err = engine.NoCache().Table("Userinfo").Cols("id, Username").Find(&tmpUsers)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(tmpUsers)

	user := &Userinfo{Uid: 1, Alias: "", Height: 0}
	affected, err := engine.Cols("Departname, Height").Update(user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println("===================", user, affected)
}

type UserCU struct {
	Id      int64
	Name    string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func testCreatedAndUpdated(engine *xorm.Engine, t *testing.T) {
	u := new(UserCU)
	err := engine.DropTables(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	u.Name = "sss"
	cnt, err := engine.Insert(u)
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

	u.Name = "xxx"
	cnt, err = engine.Id(u.Id).Update(u)
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

	u.Id = 0
	u.Created = time.Now().Add(-time.Hour * 24 * 365)
	u.Updated = u.Created
	fmt.Println(u)
	cnt, err = engine.NoAutoTime().Insert(u)
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
}

type StrangeName struct {
	Id_t int64 `xorm:"pk autoincr"`
	Name string
}

func testStrangeName(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(new(StrangeName))
	if err != nil {
		t.Error(err)
	}

	err = engine.CreateTables(new(StrangeName))
	if err != nil {
		t.Error(err)
	}

	_, err = engine.Insert(&StrangeName{Name: "sfsfdsfds"})
	if err != nil {
		t.Error(err)
	}

	beans := make([]StrangeName, 0)
	err = engine.Find(&beans)
	if err != nil {
		t.Error(err)
	}
}

func testPrefixTableName(engine *xorm.Engine, t *testing.T) {
	/*tempEngine, err := xorm.NewEngine(engine.DriverName(), engine.DataSourceName())
	//tempEngine, err := engine.Clone()
	if err != nil {
		t.Error(err)
		panic(err)
		return
	}
	defer tempEngine.Close()

	tempEngine.ShowSQL = true
	mapper := core.NewPrefixMapper(core.SnakeMapper{}, "xlw_")
	tempEngine.SetTableMapper(mapper)
	exist, err := tempEngine.IsTableExist(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
		return
	}
	if exist {
		err = tempEngine.DropTables(&Userinfo{})
		if err != nil {
			t.Error(err)
			panic(err)
			return
		}
	}

	err = tempEngine.CreateTables(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}*/
}

type CreatedUpdated struct {
	Id       int64
	Name     string
	Value    float64   `xorm:"numeric"`
	Created  time.Time `xorm:"created"`
	Created2 time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
}

func testCreatedUpdated(engine *xorm.Engine, t *testing.T) {
	err := engine.Sync(&CreatedUpdated{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	c := &CreatedUpdated{Name: "test"}
	_, err = engine.Insert(c)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	c2 := new(CreatedUpdated)
	has, err := engine.Id(c.Id).Get(c2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		panic(errors.New("no id"))
	}

	c2.Value -= 1
	_, err = engine.Id(c2.Id).Update(c2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

type Lowercase struct {
	Id    int64
	Name  string
	ended int64 `xorm:"-"`
}

func testLowerCase(engine *xorm.Engine, t *testing.T) {
	err := engine.Sync(&Lowercase{})
	_, err = engine.Where("(id) > 0").Delete(&Lowercase{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	_, err = engine.Insert(&Lowercase{ended: 1})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	ls := make([]Lowercase, 0)
	err = engine.Find(&ls)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if len(ls) != 1 {
		err = errors.New("should be 1")
		t.Error(err)
		panic(err)
	}
}

func BaseTestAll(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- directCreateTable --------------")
	directCreateTable(engine, t)
	fmt.Println("-------------- insert --------------")
	insert(engine, t)
	fmt.Println("-------------- testInsertDefault --------------")
	testInsertDefault(engine, t)
	fmt.Println("-------------- insertAutoIncr --------------")
	insertAutoIncr(engine, t)
	fmt.Println("-------------- insertMulti --------------")
	insertMulti(engine, t)
	fmt.Println("-------------- insertTwoTable --------------")
	insertTwoTable(engine, t)
	fmt.Println("-------------- testDelete --------------")
	testDelete(engine, t)
	fmt.Println("-------------- get --------------")
	get(engine, t)
	fmt.Println("-------------- testCascade --------------")
	testCascade(engine, t)
	fmt.Println("-------------- find --------------")
	find(engine, t)
	fmt.Println("-------------- find2 --------------")
	find2(engine, t)
	fmt.Println("-------------- findMap --------------")
	findMap(engine, t)
	fmt.Println("-------------- findMap2 --------------")
	findMap2(engine, t)
	fmt.Println("-------------- count --------------")
	count(engine, t)
	fmt.Println("-------------- where --------------")
	where(engine, t)
	fmt.Println("-------------- in --------------")
	in(engine, t)
	fmt.Println("-------------- limit --------------")
	limit(engine, t)
	fmt.Println("-------------- testCustomTableName --------------")
	testCustomTableName(engine, t)
	fmt.Println("-------------- testDump --------------")
	testDump(engine, t)
	fmt.Println("-------------- testConversion --------------")
	testConversion(engine, t)
}

func BaseTestAll2(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- table --------------")
	table(engine, t)
	fmt.Println("-------------- createMultiTables --------------")
	createMultiTables(engine, t)
	fmt.Println("-------------- tableOp --------------")
	tableOp(engine, t)
	fmt.Println("-------------- testCharst --------------")
	testCharst(engine, t)
	fmt.Println("-------------- testStoreEngine --------------")
	testStoreEngine(engine, t)
	fmt.Println("-------------- testExtends --------------")
	testExtends(engine, t)
	fmt.Println("-------------- testColTypes --------------")
	testColTypes(engine, t)
	fmt.Println("-------------- testCustomType --------------")
	testCustomType(engine, t)
	fmt.Println("-------------- testCreatedAndUpdated --------------")
	testCreatedAndUpdated(engine, t)
	fmt.Println("-------------- testIndexAndUnique --------------")
	testIndexAndUnique(engine, t)
	fmt.Println("-------------- testIntId --------------")
	testIntId(engine, t)
	fmt.Println("-------------- testInt32Id --------------")
	testInt32Id(engine, t)
	fmt.Println("-------------- testUintId --------------")
	testUintId(engine, t)
	fmt.Println("-------------- testUint32Id --------------")
	testUint32Id(engine, t)
	fmt.Println("-------------- testUint64Id --------------")
	testUint64Id(engine, t)
	fmt.Println("-------------- testMetaInfo --------------")
	testMetaInfo(engine, t)
	fmt.Println("-------------- testIterate --------------")
	testIterate(engine, t)
	fmt.Println("-------------- testRows --------------")
	testRows(engine, t)
	fmt.Println("-------------- testStrangeName --------------")
	testStrangeName(engine, t)
	fmt.Println("-------------- testVersion --------------")
	testVersion(engine, t)
	fmt.Println("-------------- testDistinct --------------")
	testDistinct(engine, t)
	fmt.Println("-------------- testUseBool --------------")
	testUseBool(engine, t)
	fmt.Println("-------------- testBool --------------")
	testBool(engine, t)
	fmt.Println("-------------- testTime --------------")
	testTime(engine, t)
	fmt.Println("-------------- testPrefixTableName --------------")
	testPrefixTableName(engine, t)
	fmt.Println("-------------- testCreatedUpdated --------------")
	testCreatedUpdated(engine, t)
	fmt.Println("-------------- testLowercase ---------------")
	testLowerCase(engine, t)
	fmt.Println("-------------- processors --------------")
	testProcessors(engine, t)
	fmt.Println("-------------- transaction --------------")
	transaction(engine, t)
	fmt.Println("-------------- testCacheDomain --------------")
	testCacheDomain(engine, t)
	fmt.Println("-------------- testDeleted --------------")
	testDeleted(engine, t)
	fmt.Println("-------------- testCompositeKey --------------")
	testCompositeKey(engine, t)
	fmt.Println("-------------- testCompositeKey2 --------------")
	testCompositeKey2(engine, t)
	fmt.Println("-------------- testStringPK --------------")
	testStringPK(engine, t)
}

// !nash! the 3rd set of the test is intended for non-cache enabled engine
func BaseTestAll3(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- processors TX --------------")
	testProcessorsTx(engine, t)
	fmt.Println("-------------- insert pointer data --------------")
	testPointerData(engine, t)
	fmt.Println("-------------- insert null data --------------")
	testNullValue(engine, t)
	fmt.Println("-------------- testNoCacheDomain --------------")
	testNoCacheDomain(engine, t)
}

func BaseTestAllSnakeMapper(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- query --------------")
	testQuery(engine, t)
	fmt.Println("-------------- exec --------------")
	exec(engine, t)
	fmt.Println("-------------- update --------------")
	update(engine, t)
	fmt.Println("-------------- order --------------")
	order(engine, t)
	fmt.Println("-------------- join --------------")
	join(engine, t)
	fmt.Println("-------------- having --------------")
	having(engine, t)
	fmt.Println("-------------- combineTransaction --------------")
	combineTransaction(engine, t)
	fmt.Println("-------------- testCols --------------")
	testCols(engine, t)
}

func BaseTestAllSameMapper(engine *xorm.Engine, t *testing.T) {
	fmt.Println("-------------- query --------------")
	testQuerySameMapper(engine, t)
	fmt.Println("-------------- exec --------------")
	execSameMapper(engine, t)
	fmt.Println("-------------- update --------------")
	updateSameMapper(engine, t)
	fmt.Println("-------------- order --------------")
	orderSameMapper(engine, t)
	fmt.Println("-------------- join --------------")
	joinSameMapper(engine, t)
	fmt.Println("-------------- having --------------")
	havingSameMapper(engine, t)
	fmt.Println("-------------- combineTransaction --------------")
	combineTransactionSameMapper(engine, t)
	fmt.Println("-------------- testCols --------------")
	testColsSameMapper(engine, t)
}
