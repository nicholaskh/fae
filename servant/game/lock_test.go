package game

import (
	"github.com/nicholaskh/assert"
	"github.com/nicholaskh/fae/config"
	"testing"
	"time"
)

func TestLockBasic(t *testing.T) {
	cf := &config.ConfigGame{
		LockMaxItems: 10,
		LockExpires:  10 * time.Second,
	}
	l := newLock(cf)
	k1 := "hello"
	k2 := "world"

	assert.Equal(t, true, l.Lock(k1))
	assert.Equal(t, false, l.Lock(k1))
	assert.Equal(t, true, l.Lock(k2))
	assert.Equal(t, false, l.Lock(k2))

	t.Logf("%+v", l.items)

	l.Unlock(k1)
	assert.Equal(t, true, l.Lock(k1))
	l.Unlock(k2)
	assert.Equal(t, true, l.Lock(k2))
}

func TestLockExpires(t *testing.T) {
	cf := &config.ConfigGame{
		LockMaxItems: 10,
		LockExpires:  1 * time.Second,
	}
	l := newLock(cf)
	k := "hello"
	l.Lock(k)
	assert.Equal(t, false, l.Lock(k))
	time.Sleep(time.Second)
	assert.Equal(t, true, l.Lock(k))
}

func BenchmarkLockBasic(b *testing.B) {
	cf := &config.ConfigGame{
		LockMaxItems: 10,
		LockExpires:  10 * time.Second,
	}
	l := newLock(cf)
	k := "haha"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Lock(k)
		l.Unlock(k)
	}
}
