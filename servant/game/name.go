package game

import (
	"errors"
	"math/rand"
	"time"
)

// TODO periodically dump bits to redis, when startup read from redis
const (
	NameCharMin = uint8('!') // 33, space is 32
	NameCharMax = uint8('~') // 126
)

var (
	ErrNameLen = errors.New("Name length not match slot length")
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// FIXME should use mutex to protect bitmap
type NameGen struct {
	bits     [][]byte
	DbLoaded bool
}

func newNameGen(slots int) (this *NameGen) {
	this = new(NameGen)
	this.DbLoaded = false
	this.bits = make([][]byte, slots)
	for i, _ := range this.bits {
		this.bits[i] = make([]byte, (NameCharMax-NameCharMin)/8+1)

		for char := NameCharMin; char <= NameCharMax; char++ {
			this.bits[i][this.pos(char)] = 0x00
		}
	}

	return
}

func (this *NameGen) slots() int {
	return len(this.bits)
}

func (this *NameGen) offset(char uint8) uint {
	return uint(char - NameCharMin)
}

func (this *NameGen) pos(char uint8) uint {
	return uint(this.offset(char) / 8)
}

func (this *NameGen) index(char uint8) int {
	return int(this.offset(char) % 8)
}

func (this *NameGen) getBit(slot int, char uint8) uint8 {
	index, pos := this.index(char), this.pos(char)
	if slot > len(this.bits) || char < NameCharMin || char > NameCharMax {
		return 0
	}

	return (this.bits[slot][index] >> pos) & 0x01
}

func (this *NameGen) setBit(slot int, char uint8, value uint8) {
	index, pos := this.index(char), this.pos(char)
	if slot > this.slots() || char < NameCharMin || char > NameCharMax {
		// invalid argument
		return
	}

	// value can only be 0/1
	if value == 0 {
		this.bits[slot][index] &^= 0x01 << pos
	} else {
		this.bits[slot][index] |= 0x01 << pos
	}

}

func (this *NameGen) SetBusy(name string) error {
	if len(name) != this.slots() {
		return ErrNameLen
	}

	for slot := 0; slot < this.slots(); slot++ {
		this.setBit(slot, name[slot], 1)
	}

	return nil
}

func (this *NameGen) Next() string {
	var (
		rv       string
		randChar uint8
		busyN    int
	)

	// TODO if found a busy char, find next non-busy one instead of rand
	for {
		rv = ""
		busyN = 0
		for slot := 0; slot < len(this.bits); slot++ {
			randChar = NameCharMin + uint8(rand.Int31n(int32(NameCharMax-NameCharMin)))
			if this.getBit(slot, randChar) != 0 {
				// this char in this slot is used
				busyN++
			}

			rv += string(randChar)
		}

		if busyN == this.slots() {
			// this name all bits are busy
			continue
		}

		// got it
		break
	}

	this.SetBusy(rv)

	return rv
}

func (this *NameGen) Contains(s string) bool {
	if len(s) != this.slots() {
		return true
	}

	sum := uint8(8)
	for slot := 0; slot < this.slots(); slot++ {
		sum += this.getBit(slot, s[slot])
	}

	return sum > 0
}

func (this *NameGen) Size() int {
	return len(this.bits[0]) * this.slots() * 8
}
