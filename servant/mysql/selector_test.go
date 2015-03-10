package mysql

import (
	"fmt"
	"github.com/nicholaskh/assert"
	"github.com/nicholaskh/fae/config"
	"github.com/nicholaskh/golib/server"
	"github.com/nicholaskh/golib/str"
	"strconv"
	"strings"
	"testing"
)

func TestSelectorStandardEndsWithDigit(t *testing.T) {
	s := newStandardServerSelector(new(config.ConfigMysql))
	assert.Equal(t, true, s.endsWithDigit("AllianceShard8"))
	assert.Equal(t, false, s.endsWithDigit("ShardLookup"))
}

func TestSelectorStandardPoolServers(t *testing.T) {
	s := server.NewServer("test")
	s.LoadConfig("../../etc/faed.cf.sample")
	section, _ := s.Conf.Section("servants.mysql")
	cf := &config.ConfigMysql{}
	cf.LoadConfig(section)

	sel := newStandardServerSelector(cf)

	assert.Equal(t, 1, len(sel.PoolServers("UserShard")))
}

func BenchmarkEndsWithDigit(b *testing.B) {
	s := newStandardServerSelector(new(config.ConfigMysql))
	for i := 0; i < b.N; i++ {
		s.endsWithDigit("UserShard1")
	}
}

// 143 ns/op
func BenchmarkStrconvItoa(b *testing.B) {
	b.ReportAllocs()
	pool := "UserShard"
	hintId := 12
	for i := 0; i < b.N; i++ {
		_ = pool + strconv.Itoa(hintId)
	}
}

// 307 ns/op
func BenchmarkStringsJoin(b *testing.B) {
	b.ReportAllocs()
	table := "UserLookup"
	a := []string{
		"SELECT shardId FROM",
		table,
		"WHERE entityId=?",
	}
	for i := 0; i < b.N; i++ {
		strings.Join(a, " ")
	}
}

// 184 ns/op
func BenchmarkStringBuilderConcat(b *testing.B) {
	const (
		s1 = "SELECT shardId FROM "
		s2 = " WHERE entityId=?"
	)

	b.ReportAllocs()
	table := "UserLookup"
	sb := str.NewStringBuilder()
	for i := 0; i < b.N; i++ {
		sb.WriteString(s1)
		sb.WriteString(table)
		sb.WriteString(s2)
		sb.String()
		sb.Reset()
	}
}

// 168 ns/op
func BenchmarkStringConcat(b *testing.B) {
	b.ReportAllocs()
	table := "UserLookup"
	for i := 0; i < b.N; i++ {
		_ = "SELECT shardId FROM " +
			table +
			" WHERE entityId=?"
	}
}

// 424 ns/op
func BenchmarkFmtStringConcat(b *testing.B) {
	b.ReportAllocs()
	table := "UserLookup"
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("SELECT shardId FROM %s WHERE entityId=?", table)
	}
}
