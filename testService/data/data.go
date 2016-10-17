// Package data provide test data
package data

import (
	"encoding/binary"
	"math/rand"
	"time"
)

const (
	byteL = 1 << iota
	int16L
	int32L
	int64L
)

// Reply is an immutable bytes for simple test
// value:
//   99(Uint64) + 99(Uint32) + 99(Uint16) + 99(Uint8) + 99(int64) + success!(bytes)
var Reply []byte

// init initial Reply bytes
func init() {
	t := time.Date(2016, time.October, 1, 12, 0, 0, 0, time.Local)
	replyLen := int32L + int64L + byteL + int64L + int32L + int16L + byteL + int64L + 8*byteL
	Reply = make([]byte, replyLen)
	length := 0

	binary.BigEndian.PutUint32(Reply[length:], uint32(replyLen))
	length += int32L
	binary.BigEndian.PutUint64(Reply[length:], uint64(t.UnixNano()))
	length += int64L
	Reply[length] = byte(99)
	length += byteL
	binary.BigEndian.PutUint64(Reply[length:], uint64(99))
	length += int64L
	binary.BigEndian.PutUint32(Reply[length:], uint32(99))
	length += int32L
	binary.BigEndian.PutUint16(Reply[length:], uint16(99))
	length += int16L
	Reply[length] = byte(99)
	length += byteL
	binary.BigEndian.PutUint64(Reply[length:], uint64(99))
	length += int64L

	copy(Reply[length:], "success!")
}

// RandomUserID return a random user id message.
// value:
//   randomID(uint64)
func RandomUserID() (result []byte) {
	t := time.Now()
	replyLen := int32L + int64L + byteL + int64L
	result = make([]byte, replyLen)
	length := 0

	binary.BigEndian.PutUint32(result[length:], uint32(replyLen))
	length += int32L
	binary.BigEndian.PutUint64(result[length:], uint64(t.UnixNano()))
	length += int64L
	result[length] = byte(212)
	length += byteL
	binary.BigEndian.PutUint64(result[length:], uint64(rand.Int63()))

	return
}
