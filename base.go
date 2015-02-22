package tests

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

func directCreateTable(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&Userinfo{}, &Userdetail{}, &Numeric{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Sync(&Userinfo{}, &Userdetail{}, new(Picture), new(Numeric))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	isEmpty, err := engine.IsTableEmpty(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !isEmpty {
		err = errors.New("userinfo should empty")
		t.Error(err)
		panic(err)
	}

	tbName := engine.TableMapper.Obj2Table("Userinfo")
	isEmpty, err = engine.IsTableEmpty(tbName)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !isEmpty {
		err = errors.New("userinfo should empty")
		t.Error(err)
		panic(err)
	}

	err = engine.DropTables(&Userinfo{}, &Userdetail{}, new(Numeric))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&Userinfo{}, &Userdetail{}, new(Numeric))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateIndexes(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateIndexes(&Userdetail{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateUniques(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateUniques(&Userdetail{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
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

func limit(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Limit(2, 1).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
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

func join(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Join("LEFT", "userdetail", "userinfo.id=userdetail.id").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func having(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.GroupBy("username").Having("username='xlw'").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
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

func joinSameMapper(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Join("LEFT", "`Userdetail`", "`Userinfo`.`(id)`=`Userdetail`.`Id`").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
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

func table(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables("user_user")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Table("user_user").CreateTable(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func createMultiTables(engine *xorm.Engine, t *testing.T) {
	session := engine.NewSession()
	defer session.Close()

	user := &Userinfo{}
	err := session.Begin()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for i := 0; i < 10; i++ {
		tableName := fmt.Sprintf("user_%v", i)

		err = session.DropTable(tableName)
		if err != nil {
			session.Rollback()
			t.Error(err)
			panic(err)
		}

		err = session.Table(tableName).CreateTable(user)
		if err != nil {
			session.Rollback()
			t.Error(err)
			panic(err)
		}
	}
	err = session.Commit()
	if err != nil {
		t.Error(err)
		panic(err)
	}
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

type allCols struct {
	Bit       int   `xorm:"BIT"`
	TinyInt   int8  `xorm:"TINYINT"`
	SmallInt  int16 `xorm:"SMALLINT"`
	MediumInt int32 `xorm:"MEDIUMINT"`
	Int       int   `xorm:"INT"`
	Integer   int   `xorm:"INTEGER"`
	BigInt    int64 `xorm:"BIGINT"`

	Char       string `xorm:"CHAR(12)"`
	Varchar    string `xorm:"VARCHAR(54)"`
	TinyText   string `xorm:"TINYTEXT"`
	Text       string `xorm:"TEXT"`
	MediumText string `xorm:"MEDIUMTEXT"`
	LongText   string `xorm:"LONGTEXT"`
	Binary     []byte `xorm:"BINARY(23)"`
	VarBinary  []byte `xorm:"VARBINARY(12)"`

	Date       time.Time `xorm:"DATE"`
	DateTime   time.Time `xorm:"DATETIME"`
	Time       time.Time `xorm:"TIME"`
	TimeStamp  time.Time `xorm:"TIMESTAMP"`
	TimeStampZ time.Time `xorm:"TIMESTAMPZ"`

	Decimal float64 `xorm:"DECIMAL"`
	Numeric float64 `xorm:"NUMERIC"`

	Real   float32 `xorm:"REAL"`
	Float  float32 `xorm:"FLOAT"`
	Double float64 `xorm:"DOUBLE"`

	TinyBlob   []byte `xorm:"TINYBLOB"`
	Blob       []byte `xorm:"BLOB"`
	MediumBlob []byte `xorm:"MEDIUMBLOB"`
	LongBlob   []byte `xorm:"LONGBLOB"`
	Bytea      []byte `xorm:"BYTEA"`

	Map   map[string]string `xorm:"TEXT"`
	Slice []string          `xorm:"TEXT"`

	Bool bool `xorm:"BOOL"`

	Serial int `xorm:"SERIAL"`
	//BigSerial int64 `xorm:"BIGSERIAL"`
}

func testColTypes(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(&allCols{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&allCols{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	ac := &allCols{
		1,
		4,
		8,
		16,
		32,
		64,
		128,

		"123",
		"fafdafa",
		"fafafafdsafdsafdaf",
		"fdsafafdsafdsaf",
		"fafdsafdsafdsfadasfsfafd",
		"fadfdsafdsafasfdasfds",
		[]byte("fdafsafdasfdsafsa"),
		[]byte("fdsafsdafs"),

		time.Now(),
		time.Now(),
		time.Now(),
		time.Now(),
		time.Now(),

		1.34,
		2.44302346,

		1.3344,
		2.59693523,
		3.2342523543,

		[]byte("fafdasf"),
		[]byte("fafdfdsafdsafasf"),
		[]byte("faffadsfdsdasf"),
		[]byte("faffdasfdsadasf"),
		[]byte("fafasdfsadffdasf"),

		map[string]string{"1": "1", "2": "2"},
		[]string{"1", "2", "3"},

		true,

		0,
		//21,
	}

	cnt, err := engine.Insert(ac)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert return not 1")
		t.Error(err)
		panic(err)
	}
	newAc := &allCols{}
	has, err := engine.Get(newAc)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("error no ideas")
		t.Error(err)
		panic(err)
	}

	// don't use this type as query condition
	newAc.Real = 0
	newAc.Float = 0
	newAc.Double = 0
	newAc.LongText = ""
	newAc.TinyText = ""
	newAc.MediumText = ""
	newAc.Text = ""
	newAc.Map = nil
	newAc.Slice = nil
	cnt, err = engine.Delete(newAc)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New(fmt.Sprintf("delete error, deleted counts is %v", cnt))
		t.Error(err)
		panic(err)
	}
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

type IndexOrUnique struct {
	Id        int64
	Index     int `xorm:"index"`
	Unique    int `xorm:"unique"`
	Group1    int `xorm:"index(ttt)"`
	Group2    int `xorm:"index(ttt)"`
	UniGroup1 int `xorm:"unique(lll)"`
	UniGroup2 int `xorm:"unique(lll)"`
}

func testIndexAndUnique(engine *xorm.Engine, t *testing.T) {
	err := engine.CreateTables(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.DropTables(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.CreateTables(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.CreateIndexes(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.CreateUniques(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}
}

type CacheDomain struct {
	Id   int64 `xorm:"pk cache"`
	Name string
}

func testCacheDomain(engine *xorm.Engine, t *testing.T) {
	err := engine.CreateTables(&CacheDomain{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	table := engine.TableInfo(&CacheDomain{})
	if table.Cacher == nil {
		err = fmt.Errorf("table cache is nil")
		t.Error(err)
		panic(err)
	}
}

type NoCacheDomain struct {
	Id   int64 `xorm:"pk nocache"`
	Name string
}

func testNoCacheDomain(engine *xorm.Engine, t *testing.T) {
	err := engine.CreateTables(&NoCacheDomain{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	table := engine.TableInfo(&NoCacheDomain{})
	if table.Cacher != nil {
		err = fmt.Errorf("table cache exist")
		t.Error(err)
		panic(err)
	}
}

type Deleted struct {
	Id        int64 `xorm:"pk"`
	Name      string
	DeletedAt time.Time `xorm:"deleted"`
}

func testDeleted(engine *xorm.Engine, t *testing.T) {
	err := engine.CreateTables(&Deleted{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	//table := engine.TableInfo(&Deleted{})

	engine.InsertOne(&Deleted{Id: 1, Name: "11111"})
	engine.InsertOne(&Deleted{Id: 2, Name: "22222"})
	engine.InsertOne(&Deleted{Id: 3, Name: "33333"})

	// Test normal Find()
	var records1 []Deleted
	err = engine.Where("id > 0").Find(&records1, &Deleted{})
	if len(records1) != 3 {
		t.Fatalf("Find failed: expected=%d, actual=%d, err=%v", 3, len(records1), err)
	}

	// Test normal Get()
	record1 := &Deleted{}
	has, err := engine.Id(1).Get(record1)
	if !has {
		t.Fatalf("Get failed: expected=%v, actual=%v, err=%v", true, has, err)
	}
	//fmt.Println("----- get:", record1)

	// Test Delete() with deleted
	affected, err := engine.Id(1).Delete(&Deleted{})
	if affected != 1 {
		t.Fatalf("Delete failed: expected=%v, actual=%v, err=%v", 1, affected, err)
	}
	has, err = engine.Id(1).Get(&Deleted{})
	if has {
		t.Fatalf("Delete failed. Must not get any records.")
	}
	var records2 []Deleted
	err = engine.Where("id > 0").Find(&records2)
	if len(records2) != 2 {
		t.Fatalf("Find() failed.")
	}

	// Test no rows affected after Delete() again.
	affected, err = engine.Id(1).Delete(&Deleted{})
	if affected != 0 {
		t.Fatalf("Delete failed. No rows must be affected: expected=%v, actual=%v, err=%v", 0, affected, err)
	}

	// Deleted.DeletedAt must not be updated.
	affected, err = engine.Id(2).Update(&Deleted{Name: "2", DeletedAt: time.Now()})
	if affected != 1 {
		t.Fatalf("Update failed: expected=%v, actual=%v, err=%v", 1, affected, err)
	}
	record2 := &Deleted{}
	has, err = engine.Id(2).Get(record2)
	if !record2.DeletedAt.IsZero() {
		t.Fatalf("Update failed: DeletedAt must be zero value. actual=%v", record2.DeletedAt)
	}

	// Test find all records whatever `deleted`.
	var unscopedRecords1 []Deleted
	err = engine.Unscoped().Where("id > 0").Find(&unscopedRecords1, &Deleted{})
	if len(unscopedRecords1) != 3 {
		fmt.Printf("unscopedRecords1 = %v\n", unscopedRecords1)
		t.Fatalf("Find failed: all records must be selected when engine.Unscoped()")
	}

	// Delete() must really delete a record with Unscoped()
	affected, err = engine.Unscoped().Id(1).Delete(&Deleted{})
	if affected != 1 {
		t.Fatalf("Delete failed")
	}
	var unscopedRecords2 []Deleted
	err = engine.Unscoped().Where("id > 0").Find(&unscopedRecords2, &Deleted{})
	if len(unscopedRecords2) != 2 {
		t.Fatalf("Find failed: Only 2 records must be selected when engine.Unscoped()")
	}
}

func testMetaInfo(engine *xorm.Engine, t *testing.T) {
	tables, err := engine.DBMetas()
	if err != nil {
		t.Error(err)
		panic(err)
	}

	for _, table := range tables {
		fmt.Println(table.Name)
		//for _, col := range table.Columns() {
		//TODO: engine.dialect show exported
		//fmt.Println(col.String(engine.dia))
		//}

		for _, index := range table.Indexes {
			fmt.Println(index.Name, index.Type, strings.Join(index.Cols, ","))
		}
	}
}

func testIterate(engine *xorm.Engine, t *testing.T) {
	err := engine.Omit("is_man").Iterate(new(Userinfo), func(idx int, bean interface{}) error {
		user := bean.(*Userinfo)
		fmt.Println(idx, "--", user)
		return nil
	})

	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testRows(engine *xorm.Engine, t *testing.T) {
	rows, err := engine.Omit("is_man").Rows(new(Userinfo))
	if err != nil {
		t.Error(err)
		panic(err)
	}
	defer rows.Close()

	idx := 0
	user := new(Userinfo)
	for rows.Next() {
		err = rows.Scan(user)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		fmt.Println(idx, "--", user)
		idx++
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

type VersionS struct {
	Id   int64
	Name string
	Ver  int `xorm:"version"`
}

func testVersion(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(new(VersionS))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(new(VersionS))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	ver := &VersionS{Name: "sfsfdsfds"}
	_, err = engine.Insert(ver)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(ver)
	if ver.Ver != 1 {
		err = errors.New("insert error")
		t.Error(err)
		panic(err)
	}

	newVer := new(VersionS)
	has, err := engine.Id(ver.Id).Get(newVer)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		t.Error(errors.New(fmt.Sprintf("no version id is %v", ver.Id)))
		panic(err)
	}
	fmt.Println(newVer)
	if newVer.Ver != 1 {
		err = errors.New("insert error")
		t.Error(err)
		panic(err)
	}

	newVer.Name = "-------"
	_, err = engine.Id(ver.Id).Update(newVer)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if newVer.Ver != 2 {
		err = errors.New("update should set version back to struct")
		t.Error(err)
	}

	if newVer.Ver != 2 {
		err = errors.New("update error")
		t.Error(err)
		panic(err)
	}

	newVer = new(VersionS)
	has, err = engine.Id(ver.Id).Get(newVer)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(newVer)
	if newVer.Ver != 2 {
		err = errors.New("insert error")
		t.Error(err)
		panic(err)
	}

	/*
	   newVer.Name = "-------"
	   _, err = engine.Id(ver.Id).Update(newVer)
	   if err != nil {
	       t.Error(err)
	       return
	   }*/
}

func testDistinct(engine *xorm.Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	departname := engine.TableMapper.Obj2Table("Departname")
	err := engine.Distinct(departname).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users) != 1 {
		t.Error(err)
		panic(errors.New("should be one record"))
	}

	fmt.Println(users)

	type Depart struct {
		Departname string
	}

	users2 := make([]Depart, 0)
	err = engine.Distinct(departname).Table(new(Userinfo)).Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users2) != 1 {
		t.Error(err)
		panic(errors.New("should be one record"))
	}
	fmt.Println(users2)
}

func testUseBool(engine *xorm.Engine, t *testing.T) {
	cnt1, err := engine.Count(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	users := make([]Userinfo, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	var fNumber int64
	for _, u := range users {
		if u.IsMan == false {
			fNumber += 1
		}
	}

	cnt2, err := engine.UseBool().Update(&Userinfo{IsMan: true})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if fNumber != cnt2 {
		fmt.Println("cnt1", cnt1, "fNumber", fNumber, "cnt2", cnt2)
		/*err = errors.New("Updated number is not corrected.")
		  t.Error(err)
		  panic(err)*/
	}

	_, err = engine.Update(&Userinfo{IsMan: true})
	if err == nil {
		err = errors.New("error condition")
		t.Error(err)
		panic(err)
	}
}

func testBool(engine *xorm.Engine, t *testing.T) {
	_, err := engine.UseBool().Update(&Userinfo{IsMan: true})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	users := make([]Userinfo, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		if !user.IsMan {
			err = errors.New("update bool or find bool error")
			t.Error(err)
			panic(err)
		}
	}

	_, err = engine.UseBool().Update(&Userinfo{IsMan: false})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	users = make([]Userinfo, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		if user.IsMan {
			err = errors.New("update bool or find bool error")
			t.Error(err)
			panic(err)
		}
	}
}

type TTime struct {
	Id int64
	T  time.Time
	Tz time.Time `xorm:"timestampz"`
}

func (t *TTime) String() string {
	return fmt.Sprintf("%v|T:%v|Tz:%v", t.Id, t.T, t.Tz)
}

func testTime(engine *xorm.Engine, t *testing.T) {
	err := engine.Sync(&TTime{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tt := &TTime{}

	println("b4 Insert tt:", tt.String())
	_, err = engine.Insert(tt)

	println("after Insert tt:", tt.String())
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tt2 := &TTime{Id: tt.Id}
	println("b4 Get tt2:", tt2.String())
	has, err := engine.Get(tt2)
	println("after Get tt2:", tt2.String())
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("no record error")
		t.Error(err)
		panic(err)
	}

	tt3 := &TTime{T: time.Now(), Tz: time.Now()}
	println("b4 Insert tt3:", tt3.String())
	_, err = engine.Insert(tt3)
	println("after Insert tt3:", tt3.String())
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tt4s := make([]TTime, 0)
	println("b4 Insert tt4s:", tt4s)
	err = engine.Find(&tt4s)
	println("after Insert tt4s:", tt4s)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println("=======\n", tt4s, "=======\n")
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

type NullData struct {
	Id         int64
	StringPtr  *string
	StringPtr2 *string `xorm:"text"`
	BoolPtr    *bool
	BytePtr    *byte
	UintPtr    *uint
	Uint8Ptr   *uint8
	Uint16Ptr  *uint16
	Uint32Ptr  *uint32
	Uint64Ptr  *uint64
	IntPtr     *int
	Int8Ptr    *int8
	Int16Ptr   *int16
	Int32Ptr   *int32
	Int64Ptr   *int64
	RunePtr    *rune
	Float32Ptr *float32
	Float64Ptr *float64
	// Complex64Ptr *complex64 // !nashtsai! XORM yet support complex128:  'json: unsupported type: complex128'
	// Complex128Ptr *complex128 // !nashtsai! XORM yet support complex128:  'json: unsupported type: complex128'
	TimePtr *time.Time
}

type NullData2 struct {
	Id         int64
	StringPtr  string
	StringPtr2 string `xorm:"text"`
	BoolPtr    bool
	BytePtr    byte
	UintPtr    uint
	Uint8Ptr   uint8
	Uint16Ptr  uint16
	Uint32Ptr  uint32
	Uint64Ptr  uint64
	IntPtr     int
	Int8Ptr    int8
	Int16Ptr   int16
	Int32Ptr   int32
	Int64Ptr   int64
	RunePtr    rune
	Float32Ptr float32
	Float64Ptr float64
	// Complex64Ptr complex64 // !nashtsai! XORM yet support complex128:  'json: unsupported type: complex128'
	// Complex128Ptr complex128 // !nashtsai! XORM yet support complex128:  'json: unsupported type: complex128'
	TimePtr time.Time
}

type NullData3 struct {
	Id        int64
	StringPtr *string
}

func testPointerData(engine *xorm.Engine, t *testing.T) {

	err := engine.DropTables(&NullData{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&NullData{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	nullData := NullData{
		StringPtr:  new(string),
		StringPtr2: new(string),
		BoolPtr:    new(bool),
		BytePtr:    new(byte),
		UintPtr:    new(uint),
		Uint8Ptr:   new(uint8),
		Uint16Ptr:  new(uint16),
		Uint32Ptr:  new(uint32),
		Uint64Ptr:  new(uint64),
		IntPtr:     new(int),
		Int8Ptr:    new(int8),
		Int16Ptr:   new(int16),
		Int32Ptr:   new(int32),
		Int64Ptr:   new(int64),
		RunePtr:    new(rune),
		Float32Ptr: new(float32),
		Float64Ptr: new(float64),
		// Complex64Ptr: new(complex64),
		// Complex128Ptr: new(complex128),
		TimePtr: new(time.Time),
	}

	*nullData.StringPtr = "abc"
	*nullData.StringPtr2 = "123"
	*nullData.BoolPtr = true
	*nullData.BytePtr = 1
	*nullData.UintPtr = 1
	*nullData.Uint8Ptr = 1
	*nullData.Uint16Ptr = 1
	*nullData.Uint32Ptr = 1
	*nullData.Uint64Ptr = 1
	*nullData.IntPtr = -1
	*nullData.Int8Ptr = -1
	*nullData.Int16Ptr = -1
	*nullData.Int32Ptr = -1
	*nullData.Int64Ptr = -1
	*nullData.RunePtr = 1
	*nullData.Float32Ptr = -1.2
	*nullData.Float64Ptr = -1.1
	// *nullData.Complex64Ptr = 123456789012345678901234567890
	// *nullData.Complex128Ptr = 123456789012345678901234567890123456789012345678901234567890
	*nullData.TimePtr = time.Now()

	cnt, err := engine.Insert(&nullData)
	fmt.Println(nullData.Id)
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
	if nullData.Id <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	// verify get values
	nullDataGet := NullData{}
	has, err := engine.Id(nullData.Id).Get(&nullDataGet)
	if err != nil {
		t.Error(err)
		panic(err)
	} else if !has {
		t.Error(errors.New("ID not found"))
	}

	if *nullDataGet.StringPtr != *nullData.StringPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.StringPtr)))
	}

	if *nullDataGet.StringPtr2 != *nullData.StringPtr2 {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.StringPtr2)))
	}

	if *nullDataGet.BoolPtr != *nullData.BoolPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%t]", *nullDataGet.BoolPtr)))
	}

	if *nullDataGet.UintPtr != *nullData.UintPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.UintPtr)))
	}

	if *nullDataGet.Uint8Ptr != *nullData.Uint8Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Uint8Ptr)))
	}

	if *nullDataGet.Uint16Ptr != *nullData.Uint16Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Uint16Ptr)))
	}

	if *nullDataGet.Uint32Ptr != *nullData.Uint32Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Uint32Ptr)))
	}

	if *nullDataGet.Uint64Ptr != *nullData.Uint64Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Uint64Ptr)))
	}

	if *nullDataGet.IntPtr != *nullData.IntPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.IntPtr)))
	}

	if *nullDataGet.Int8Ptr != *nullData.Int8Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Int8Ptr)))
	}

	if *nullDataGet.Int16Ptr != *nullData.Int16Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Int16Ptr)))
	}

	if *nullDataGet.Int32Ptr != *nullData.Int32Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Int32Ptr)))
	}

	if *nullDataGet.Int64Ptr != *nullData.Int64Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Int64Ptr)))
	}

	if *nullDataGet.RunePtr != *nullData.RunePtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.RunePtr)))
	}

	if *nullDataGet.Float32Ptr != *nullData.Float32Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Float32Ptr)))
	}

	if *nullDataGet.Float64Ptr != *nullData.Float64Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Float64Ptr)))
	}

	// if *nullDataGet.Complex64Ptr != *nullData.Complex64Ptr {
	//  t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Complex64Ptr)))
	// }

	// if *nullDataGet.Complex128Ptr != *nullData.Complex128Ptr {
	//  t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Complex128Ptr)))
	// }

	/*if (*nullDataGet.TimePtr).Unix() != (*nullData.TimePtr).Unix() {
	      t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]:[%v]", *nullDataGet.TimePtr, *nullData.TimePtr)))
	  } else {
	      // !nashtsai! mymysql driver will failed this test case, due the time is roundup to nearest second, I would considered this is a bug in mymysql driver
	      fmt.Printf("time value: [%v]:[%v]", *nullDataGet.TimePtr, *nullData.TimePtr)
	      fmt.Println()
	  }*/
	// --

	// using instance type should just work too
	nullData2Get := NullData2{}

	tableName := engine.TableMapper.Obj2Table("NullData")

	has, err = engine.Table(tableName).Id(nullData.Id).Get(&nullData2Get)
	if err != nil {
		t.Error(err)
		panic(err)
	} else if !has {
		t.Error(errors.New("ID not found"))
	}

	if nullData2Get.StringPtr != *nullData.StringPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.StringPtr)))
	}

	if nullData2Get.StringPtr2 != *nullData.StringPtr2 {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.StringPtr2)))
	}

	if nullData2Get.BoolPtr != *nullData.BoolPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%t]", nullData2Get.BoolPtr)))
	}

	if nullData2Get.UintPtr != *nullData.UintPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.UintPtr)))
	}

	if nullData2Get.Uint8Ptr != *nullData.Uint8Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Uint8Ptr)))
	}

	if nullData2Get.Uint16Ptr != *nullData.Uint16Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Uint16Ptr)))
	}

	if nullData2Get.Uint32Ptr != *nullData.Uint32Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Uint32Ptr)))
	}

	if nullData2Get.Uint64Ptr != *nullData.Uint64Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Uint64Ptr)))
	}

	if nullData2Get.IntPtr != *nullData.IntPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.IntPtr)))
	}

	if nullData2Get.Int8Ptr != *nullData.Int8Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Int8Ptr)))
	}

	if nullData2Get.Int16Ptr != *nullData.Int16Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Int16Ptr)))
	}

	if nullData2Get.Int32Ptr != *nullData.Int32Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Int32Ptr)))
	}

	if nullData2Get.Int64Ptr != *nullData.Int64Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Int64Ptr)))
	}

	if nullData2Get.RunePtr != *nullData.RunePtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.RunePtr)))
	}

	if nullData2Get.Float32Ptr != *nullData.Float32Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Float32Ptr)))
	}

	if nullData2Get.Float64Ptr != *nullData.Float64Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Float64Ptr)))
	}

	// if nullData2Get.Complex64Ptr != *nullData.Complex64Ptr {
	//  t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Complex64Ptr)))
	// }

	// if nullData2Get.Complex128Ptr != *nullData.Complex128Ptr {
	//  t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", nullData2Get.Complex128Ptr)))
	// }

	/*if nullData2Get.TimePtr.Unix() != (*nullData.TimePtr).Unix() {
	      t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]:[%v]", nullData2Get.TimePtr, *nullData.TimePtr)))
	  } else {
	      // !nashtsai! mymysql driver will failed this test case, due the time is roundup to nearest second, I would considered this is a bug in mymysql driver
	      fmt.Printf("time value: [%v]:[%v]", nullData2Get.TimePtr, *nullData.TimePtr)
	      fmt.Println()
	  }*/
	// --
}

func testNullValue(engine *xorm.Engine, t *testing.T) {

	err := engine.DropTables(&NullData{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&NullData{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	nullData := NullData{}

	cnt, err := engine.Insert(&nullData)
	fmt.Println(nullData.Id)
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
	if nullData.Id <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	nullDataGet := NullData{}

	has, err := engine.Id(nullData.Id).Get(&nullDataGet)
	if err != nil {
		t.Error(err)
		panic(err)
	} else if !has {
		t.Error(errors.New("ID not found"))
	}

	if nullDataGet.StringPtr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.StringPtr)))
	}

	if nullDataGet.StringPtr2 != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.StringPtr2)))
	}

	if nullDataGet.BoolPtr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%t]", *nullDataGet.BoolPtr)))
	}

	if nullDataGet.UintPtr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.UintPtr)))
	}

	if nullDataGet.Uint8Ptr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Uint8Ptr)))
	}

	if nullDataGet.Uint16Ptr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Uint16Ptr)))
	}

	if nullDataGet.Uint32Ptr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Uint32Ptr)))
	}

	if nullDataGet.Uint64Ptr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Uint64Ptr)))
	}

	if nullDataGet.IntPtr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.IntPtr)))
	}

	if nullDataGet.Int8Ptr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Int8Ptr)))
	}

	if nullDataGet.Int16Ptr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Int16Ptr)))
	}

	if nullDataGet.Int32Ptr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Int32Ptr)))
	}

	if nullDataGet.Int64Ptr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Int64Ptr)))
	}

	if nullDataGet.RunePtr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.RunePtr)))
	}

	if nullDataGet.Float32Ptr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Float32Ptr)))
	}

	if nullDataGet.Float64Ptr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Float64Ptr)))
	}

	// if nullDataGet.Complex64Ptr != nil {
	//  t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Complex64Ptr)))
	// }

	// if nullDataGet.Complex128Ptr != nil {
	//  t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Complex128Ptr)))
	// }

	if nullDataGet.TimePtr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.TimePtr)))
	}

	nullDataUpdate := NullData{
		StringPtr:  new(string),
		StringPtr2: new(string),
		BoolPtr:    new(bool),
		BytePtr:    new(byte),
		UintPtr:    new(uint),
		Uint8Ptr:   new(uint8),
		Uint16Ptr:  new(uint16),
		Uint32Ptr:  new(uint32),
		Uint64Ptr:  new(uint64),
		IntPtr:     new(int),
		Int8Ptr:    new(int8),
		Int16Ptr:   new(int16),
		Int32Ptr:   new(int32),
		Int64Ptr:   new(int64),
		RunePtr:    new(rune),
		Float32Ptr: new(float32),
		Float64Ptr: new(float64),
		// Complex64Ptr: new(complex64),
		// Complex128Ptr: new(complex128),
		TimePtr: new(time.Time),
	}

	*nullDataUpdate.StringPtr = "abc"
	*nullDataUpdate.StringPtr2 = "123"
	*nullDataUpdate.BoolPtr = true
	*nullDataUpdate.BytePtr = 1
	*nullDataUpdate.UintPtr = 1
	*nullDataUpdate.Uint8Ptr = 1
	*nullDataUpdate.Uint16Ptr = 1
	*nullDataUpdate.Uint32Ptr = 1
	*nullDataUpdate.Uint64Ptr = 1
	*nullDataUpdate.IntPtr = -1
	*nullDataUpdate.Int8Ptr = -1
	*nullDataUpdate.Int16Ptr = -1
	*nullDataUpdate.Int32Ptr = -1
	*nullDataUpdate.Int64Ptr = -1
	*nullDataUpdate.RunePtr = 1
	*nullDataUpdate.Float32Ptr = -1.2
	*nullDataUpdate.Float64Ptr = -1.1
	// *nullDataUpdate.Complex64Ptr = 123456789012345678901234567890
	// *nullDataUpdate.Complex128Ptr = 123456789012345678901234567890123456789012345678901234567890
	*nullDataUpdate.TimePtr = time.Now()

	cnt, err = engine.Id(nullData.Id).Update(&nullDataUpdate)
	if err != nil {
		t.Error(err)
		panic(err)
	} else if cnt != 1 {
		t.Error(errors.New("update count == 0, how can this happen!?"))
		return
	}

	// verify get values
	nullDataGet = NullData{}
	has, err = engine.Id(nullData.Id).Get(&nullDataGet)
	if err != nil {
		t.Error(err)
		return
	} else if !has {
		t.Error(errors.New("ID not found"))
		return
	}

	if *nullDataGet.StringPtr != *nullDataUpdate.StringPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.StringPtr)))
	}

	if *nullDataGet.StringPtr2 != *nullDataUpdate.StringPtr2 {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.StringPtr2)))
	}

	if *nullDataGet.BoolPtr != *nullDataUpdate.BoolPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%t]", *nullDataGet.BoolPtr)))
	}

	if *nullDataGet.UintPtr != *nullDataUpdate.UintPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.UintPtr)))
	}

	if *nullDataGet.Uint8Ptr != *nullDataUpdate.Uint8Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Uint8Ptr)))
	}

	if *nullDataGet.Uint16Ptr != *nullDataUpdate.Uint16Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Uint16Ptr)))
	}

	if *nullDataGet.Uint32Ptr != *nullDataUpdate.Uint32Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Uint32Ptr)))
	}

	if *nullDataGet.Uint64Ptr != *nullDataUpdate.Uint64Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Uint64Ptr)))
	}

	if *nullDataGet.IntPtr != *nullDataUpdate.IntPtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.IntPtr)))
	}

	if *nullDataGet.Int8Ptr != *nullDataUpdate.Int8Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Int8Ptr)))
	}

	if *nullDataGet.Int16Ptr != *nullDataUpdate.Int16Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Int16Ptr)))
	}

	if *nullDataGet.Int32Ptr != *nullDataUpdate.Int32Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Int32Ptr)))
	}

	if *nullDataGet.Int64Ptr != *nullDataUpdate.Int64Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Int64Ptr)))
	}

	if *nullDataGet.RunePtr != *nullDataUpdate.RunePtr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.RunePtr)))
	}

	if *nullDataGet.Float32Ptr != *nullDataUpdate.Float32Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Float32Ptr)))
	}

	if *nullDataGet.Float64Ptr != *nullDataUpdate.Float64Ptr {
		t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Float64Ptr)))
	}

	// if *nullDataGet.Complex64Ptr != *nullDataUpdate.Complex64Ptr {
	//  t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Complex64Ptr)))
	// }

	// if *nullDataGet.Complex128Ptr != *nullDataUpdate.Complex128Ptr {
	//  t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]", *nullDataGet.Complex128Ptr)))
	// }

	// !nashtsai! skipped mymysql test due to driver will round up time caused inaccuracy comparison
	// skipped postgres test due to postgres driver doesn't read time.Time's timzezone info when stored in the db
	// mysql and sqlite3 seem have done this correctly by storing datatime in UTC timezone, I think postgres driver
	// prefer using timestamp with timezone to sovle the issue
	if engine.DriverName() != core.POSTGRES && engine.DriverName() != "mymysql" &&
		engine.DriverName() != core.MYSQL {
		if (*nullDataGet.TimePtr).Unix() != (*nullDataUpdate.TimePtr).Unix() {
			t.Error(errors.New(fmt.Sprintf("inserted value unmatch: [%v]:[%v]", *nullDataGet.TimePtr, *nullDataUpdate.TimePtr)))
		} else {
			// !nashtsai! mymysql driver will failed this test case, due the time is roundup to nearest second, I would considered this is a bug in mymysql driver
			//  inserted value unmatch: [2013-12-25 12:12:45 +0800 CST]:[2013-12-25 12:12:44.878903653 +0800 CST]
			fmt.Printf("time value: [%v]:[%v]", *nullDataGet.TimePtr, *nullDataUpdate.TimePtr)
			fmt.Println()
		}
	}

	// update to null values
	nullDataUpdate = NullData{}

	string_ptr := engine.ColumnMapper.Obj2Table("StringPtr")

	cnt, err = engine.Id(nullData.Id).Cols(string_ptr).Update(&nullDataUpdate)
	if err != nil {
		t.Error(err)
		panic(err)
	} else if cnt != 1 {
		t.Error(errors.New("update count == 0, how can this happen!?"))
		return
	}

	// verify get values
	nullDataGet = NullData{}
	has, err = engine.Id(nullData.Id).Get(&nullDataGet)
	if err != nil {
		t.Error(err)
		return
	} else if !has {
		t.Error(errors.New("ID not found"))
		return
	}

	fmt.Printf("%+v", nullDataGet)
	fmt.Println()

	if nullDataGet.StringPtr != nil {
		t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.StringPtr)))
	}
	/*
	  if nullDataGet.StringPtr2 != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.StringPtr2)))
	  }

	  if nullDataGet.BoolPtr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%t]", *nullDataGet.BoolPtr)))
	  }

	  if nullDataGet.UintPtr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.UintPtr)))
	  }

	  if nullDataGet.Uint8Ptr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Uint8Ptr)))
	  }

	  if nullDataGet.Uint16Ptr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Uint16Ptr)))
	  }

	  if nullDataGet.Uint32Ptr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Uint32Ptr)))
	  }

	  if nullDataGet.Uint64Ptr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Uint64Ptr)))
	  }

	  if nullDataGet.IntPtr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.IntPtr)))
	  }

	  if nullDataGet.Int8Ptr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Int8Ptr)))
	  }

	  if nullDataGet.Int16Ptr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Int16Ptr)))
	  }

	  if nullDataGet.Int32Ptr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Int32Ptr)))
	  }

	  if nullDataGet.Int64Ptr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Int64Ptr)))
	  }

	  if nullDataGet.RunePtr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.RunePtr)))
	  }

	  if nullDataGet.Float32Ptr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Float32Ptr)))
	  }

	  if nullDataGet.Float64Ptr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Float64Ptr)))
	  }

	  // if nullDataGet.Complex64Ptr != nil {
	  //  t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Float64Ptr)))
	  // }

	  // if nullDataGet.Complex128Ptr != nil {
	  //  t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.Float64Ptr)))
	  // }

	  if nullDataGet.TimePtr != nil {
	      t.Error(errors.New(fmt.Sprintf("not null value: [%v]", *nullDataGet.TimePtr)))
	  }*/
	// --

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

type CustomTableName struct {
	Id   int64
	Name string
}

func (c *CustomTableName) TableName() string {
	return "customtablename"
}

func testCustomTableName(engine *xorm.Engine, t *testing.T) {
	c := new(CustomTableName)
	err := engine.DropTables(c)
	if err != nil {
		t.Error(err)
	}

	err = engine.CreateTables(c)
	if err != nil {
		t.Error(err)
	}
}

func testDump(engine *xorm.Engine, t *testing.T) {
	fp := engine.Dialect().URI().DbName + ".sql"
	os.Remove(fp)
	err := engine.DumpAllToFile(fp)
	if err != nil {
		t.Error(err)
		fmt.Println(err)
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
