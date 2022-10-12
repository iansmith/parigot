package parigot

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano())) // reproducibility
}

type ServiceId int64

func NewServiceId() ServiceId {
	id := newId(0x73) //'s'
	sid := ServiceId(id)
	return sid
}

func AsIdShort(n int64) string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(n))
	switch buf[7] {
	case 's':
		buf[7] = 0
		buf[6] = 0
		buf[5] = 0
		buf[4] = 0
		buf[3] = 0
		buf[2] = 0
		u := binary.LittleEndian.Uint64(buf)
		return fmt.Sprintf("s-%04x", u)
	}
	panic("unable to understand id:" + fmt.Sprintf("%v", buf))
}
func AsId(n int64) string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(n))
	switch buf[7] {
	case 's':
		buf[7] = 0
		u := binary.LittleEndian.Uint64(buf)
		return fmt.Sprintf("s-%x", u)
	}
	panic("unable to understand id:" + fmt.Sprintf("%v", buf))
}

func newId(s byte) int64 {
	u := rand.Uint64()
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(u))
	buf[7] = s
	x := binary.LittleEndian.Uint64(buf)
	return int64(x)
}

func (s ServiceId) String() string {
	return AsIdShort(int64(s))
}
func (s ServiceId) StringFull() string {
	return AsId(int64(s))
}
