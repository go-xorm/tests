package tests

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

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
