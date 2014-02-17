package servant

import (
	"github.com/funkygao/fae/servant/mongo"
	log "github.com/funkygao/log4go"
	"labix.org/v2/mgo/bson"
)

func (this *FunServantImpl) mongoSession(pool string,
	shardId int32) (*mongo.Session, error) {
	sess, err := this.mg.Session(pool, shardId)
	if err != nil {
		log.Error("{pool^%s id^%d} %s", pool, shardId, err)
		return nil, err
	}

	return sess, err
}

// specs: inbound params use json
func (this *FunServantImpl) unmarshalIn(d []byte) (v bson.M, err error) {
	err = bson.Unmarshal(d, &v)
	if err != nil {
		log.Error("unmarshalIn error: %s -> %s", d, err)
	}

	return
}

// specs: outbound data use bson
func (this *FunServantImpl) marshalOut(d bson.M) []byte {
	val, err := bson.Marshal(d)
	if err != nil {
		// should never happen
		log.Critical("marshalOut error: %+v -> %v", d, err)
	}
	return val
}

func (this *FunServantImpl) mgFieldsIsNil(fields []byte) bool {
	return len(fields) <= 5
}