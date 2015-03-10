package servant

import (
	"encoding/json"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/nicholaskh/fae/config"
	"github.com/nicholaskh/fae/servant/gen-go/fun/rpc"
	"github.com/nicholaskh/fae/servant/proxy"
	"github.com/nicholaskh/golib/server"
	conf "github.com/nicholaskh/jsconf"
	"github.com/nicholaskh/msgpack"
	"labix.org/v2/mgo/bson"
	"strings"
	"testing"
)

func setupServant() *FunServantImpl {
	cf, _ := conf.Load("../etc/faed.cf")
	config.LoadEngineConfig("../etc/faed.cf", cf)
	server.LaunchHttpServ(":9999", "")
	return NewFunServant(config.Engine.Servants)
}

// 683 ns/op
func BenchmarkSprintfSql(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("SELECT %s FROM %s WHERE %s", "*", "UserInfo", "uid=?")
	}

}

// 166 ns/op
func BenchmarkRawConcatSql(b *testing.B) {
	b.ReportAllocs()
	table := "UserInfo"
	column := "*"
	where := "uid=?"
	for i := 0; i < b.N; i++ {
		_ = "SELECT " + column + " FROM " + table + " WHERE " + where
	}
}

// 103 ns/op
func BenchmarkDefer(b *testing.B) {
	b.ReportAllocs()
	f := func() {
		defer func() {

		}()
	}
	for i := 0; i < b.N; i++ {
		f()
	}
}

func BenchmarkMcSet(b *testing.B) {
	servant := setupServant()
	b.ReportAllocs()

	ctx := rpc.NewContext()
	ctx.Reason = "benchmark"
	ctx.Rid = "1"
	data := rpc.NewTMemcacheData()
	data.Data = []byte("bar")
	for i := 0; i < b.N; i++ {
		servant.McSet(ctx, "default", "foo", data, 0)
	}
}

func BenchmarkGoMap(b *testing.B) {
	b.ReportAllocs()
	var x map[int]bool = make(map[int]bool)
	for i := 0; i < b.N; i++ {
		x[i] = true
	}
}

func BenchmarkJson(b *testing.B) {
	m := bson.M{
		"userId": 343434,
		"gendar": "F",
		"info": bson.M{
			"city":    "beijing",
			"hobbies": []string{"a", "b"}}}
	for i := 0; i < b.N; i++ {
		json.Marshal(m)
	}
}

func BenchmarkBson(b *testing.B) {
	m := bson.M{
		"userId": 343434,
		"gendar": "F",
		"info": bson.M{
			"city":    "beijing",
			"hobbies": []string{"a", "b"}}}
	for i := 0; i < b.N; i++ {
		bson.Marshal(m)
	}
}

// 880 ns/op
func BenchmarkIsSelectQuery(b *testing.B) {
	b.ReportAllocs()
	var sql = "select * from UserInfo where uid=? and power>?"
	for i := 0; i < b.N; i++ {
		strings.HasPrefix(strings.ToLower(sql), "select")
	}
}

// 9 ns/op
func BenchmarkIsSelectQueryWithoutLowcase(b *testing.B) {
	b.ReportAllocs()
	var sql = "select * from UserInfo where uid=? and power>?"
	for i := 0; i < b.N; i++ {
		strings.HasPrefix(sql, "select")
	}
}

func BenchmarkMsgPackSerialize(b *testing.B) {
	b.ReportAllocs()

	mysqlResult := rpc.NewMysqlResult()
	mysqlResult.Cols = make([]string, 0)
	mysqlResult.Rows = make([][]string, 0)
	colsN, rowsN := 5, 10
	for i := 0; i < colsN; i++ {
		mysqlResult.Cols = append(mysqlResult.Cols, "username")
	}
	for i := 0; i < rowsN; i++ {
		row := make([]string, 0)
		for j := 0; j < colsN; j++ {
			row = append(row, "beijing, los angels")
		}

		mysqlResult.Rows = append(mysqlResult.Rows, row)
	}

	for i := 0; i < b.N; i++ {
		msgpack.Marshal(mysqlResult)
	}

	b.SetBytes(int64(len(mysqlResult.String())))
}

func BenchmarkThriftSerialize(b *testing.B) {
	b.ReportAllocs()

	transport := thrift.NewTMemoryBuffer()
	protocol := thrift.NewTBinaryProtocolFactoryDefault()
	//iprot := protocol.GetProtocol(transport)
	oprot := protocol.GetProtocol(transport)
	mysqlResult := rpc.NewMysqlResult()
	mysqlResult.Cols = make([]string, 0)
	mysqlResult.Rows = make([][]string, 0)
	colsN, rowsN := 5, 10
	for i := 0; i < colsN; i++ {
		mysqlResult.Cols = append(mysqlResult.Cols, "username")
	}
	for i := 0; i < rowsN; i++ {
		row := make([]string, 0)
		for j := 0; j < colsN; j++ {
			row = append(row, "beijing, los angels")
		}

		mysqlResult.Rows = append(mysqlResult.Rows, row)
	}
	b.Logf("%s", mysqlResult.String())

	for i := 0; i < b.N; i++ {
		mysqlResult.Write(oprot)
	}

	transport.Close()
	b.SetBytes(int64(len(mysqlResult.String())))
}

func BenchmarkPingOnLocalhost(b *testing.B) {
	b.ReportAllocs()

	cf := &config.ConfigProxy{PoolCapacity: 1}
	client, err := proxy.New(cf).ServantByAddr("localhost:9001")
	if err != nil {
		b.Fatal(err)
	}
	defer client.Recycle()

	ctx := rpc.NewContext()
	ctx.Reason = "benchmark"
	ctx.Rid = "1"

	for i := 0; i < b.N; i++ {
		client.Ping(ctx)
	}
}
