package uuid

import (
	"encoding/hex"
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().Unix()))

// UUID represents uuid
type UUID [16]byte

// New creates new random UUID
func New() UUID {
	var uuid UUID
	_, err := random.Read(uuid[:])
	if err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid
}

func (uuid UUID) String() string {
	var res [36]byte
	hex.Encode(res[:], uuid[:4])
	hex.Encode(res[9:], uuid[4:6])
	hex.Encode(res[14:], uuid[6:8])
	hex.Encode(res[19:], uuid[8:10])
	hex.Encode(res[24:], uuid[10:])
	res[8] = '-'
	res[13] = '-'
	res[18] = '-'
	res[23] = '-'
	return string(res[:])
}

var zero UUID

func (uuid UUID) IsZero() bool {
	return uuid == zero
}
