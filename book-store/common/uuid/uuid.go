package uuid

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/Vesninovich/go-tasks/book-store/common/commonerrors"
)

// REGEX is simple regexp for UUID string
const REGEX = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"

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

// FromBytes returns UUID from given bytes slice
func FromBytes(bytes []byte) (UUID, error) {
	if len(bytes) != 16 {
		return zero, &commonerrors.InvalidInput{Reason: "length of bytes slice for uuid must be 16"}
	}
	var res UUID
	for i, b := range bytes {
		res[i] = b
	}
	// TODO: add check for correctness
	return res, nil
}

// FromString returns UUID from given UUID string
func FromString(str string) (res UUID, err error) {
	if len(str) != 36 || str[8] != '-' || str[13] != '-' || str[18] != '-' || str[23] != '-' {
		return zero, fmt.Errorf("Malformed UUID string: %s", str)
	}
	b := []byte(str)
	_, err = hex.Decode(res[:], b[:8])
	if err != nil {
		return
	}
	_, err = hex.Decode(res[4:], b[9:13])
	if err != nil {
		return
	}
	_, err = hex.Decode(res[6:], b[14:18])
	if err != nil {
		return
	}
	_, err = hex.Decode(res[8:], b[19:23])
	if err != nil {
		return
	}
	_, err = hex.Decode(res[10:], b[24:])
	return
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

// IsZero checks if UUID is zero value
func (uuid UUID) IsZero() bool {
	return uuid == zero
}
