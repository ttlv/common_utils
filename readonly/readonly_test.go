package readonly_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/ttlv/common_utils/readonly"
	"github.com/ttlv/common_utils/testingtool"
	"github.com/ttlv/common_utils/utils"
)

var (
	Admin  *admin.Admin
	Server *httptest.Server
	DB     *gorm.DB
)

type Trade struct {
	gorm.Model
	Market     string
	Price      float64
	ExchangeID uint
}

type Exchange struct {
	gorm.Model
	Name string
}

type Order struct {
	gorm.Model
	Price float64
}

func setup(max ...interface{}) {
	var err error
	DB, err = gorm.Open("mysql", "root:abc123@tcp(127.0.0.1:9300)/common_test?parseTime=True&loc=UTC&charset=utf8")
	//DB.LogMode(true)
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&Trade{}, &Order{}, &Exchange{})

	if os.Getenv("KEEP_DB") != "true" {
		utils.RunSQL(DB, `
		  TRUNCATE TABLE trades;
          TRUNCATE TABLE exchanges;
          INSERT exchanges (name) VALUES ('测试1');
          INSERT exchanges (name) VALUES ('测试2');
		  INSERT trades (market, price, exchange_id) VALUES ('btceth', 1000, 1)
		  INSERT trades (market, price, exchange_id) VALUES ('ethocx', 2000, 2)
		`)
	}
	/*for i := 0; i < 9900; i++ {
		utils.RunSQL(DB, `
		  INSERT trades (market, price) VALUES ('btceth', 1000)
		`)
	}*/

	Admin = admin.New(&admin.AdminConfig{DB: DB})
	res := Admin.AddResource(&Trade{})
	res.Meta(&admin.Meta{
		Name: "ExchangeID",
		FormattedValuer: func(record interface{}, context *qor.Context) (result interface{}) {
			exchange := Exchange{}
			DB.Where("id = ?", record.(*Trade).ExchangeID).Find(&exchange)
			return exchange.Name
		},
	})
	Admin.AddResource(&Order{})
	res.Filter(&admin.Filter{Name: "market"})
	res.UseTheme("readonly")
	if len(max) > 0 {
		readonly.Setup(Admin, max...)
	} else {
		readonly.Setup(Admin)
	}
	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)
	testingtool.TestServerDomain = "http://localhost:5000/"
	if os.Getenv("MODE") == "server" {
		fmt.Printf("URL: %v\n", "http://localhost:5000/admin")
		if err := http.ListenAndServe(":5000", mux); err != nil {
			panic(err)
		}
	} else {
		Server = httptest.NewServer(mux)
		testingtool.TestServerDomain = Server.URL
	}
}

func TestExport(t *testing.T) {
	setup()
	body, _ := testingtool.Get("/admin/trades/export?timestamp=1545982369663")
	expect := `
        ID,Market,Price,Exchange ID,Market
        2,ethocx,2000,测试2,ethocx
        1,btceth,1000,测试1,btceth
	`
	testingtool.CompareRecords(t, expect, strings.Replace(body, "\n", ";", -1))
}

func TestExportMax(t *testing.T) {
	setup(1)
	body, _ := testingtool.Get("/admin/trades/export?timestamp=1545982369663")
	expect := `下载结果超过1, 请过滤数据再下载.`
	testingtool.CompareRecords(t, expect, strings.Replace(body, "\n", ";", -1))

	body, _ = testingtool.Get("/admin/trades/export?timestamp=1545982369663&max=10")
	expect = `
	  ID,Market,Price,Exchange ID,Market
	  2,ethocx,2000,测试2,ethocx
	  1,btceth,1000,测试1,btceth
	`
	testingtool.CompareRecords(t, expect, strings.Replace(body, "\n", ";", -1))

	body, _ = testingtool.Get("/admin/trades/export?timestamp=1545982369663&max=1000000000")
	expect = `下载结果超过1, 请过滤数据再下载.`
	testingtool.CompareRecords(t, expect, strings.Replace(body, "\n", ";", -1))
}
