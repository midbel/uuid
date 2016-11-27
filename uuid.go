//Package uuid provides a basic implementation of Universal Unique Identifier as
//described in RFC 4122.
package uuid

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

const variant = 0x40

const devRandom = "/dev/urandom"

const pattern = "%08x-%04x-%04x-%02x%02x-%x"

var (
	timeMU sync.Mutex
	randMU sync.Mutex
)

var (
	DNS = &UUID{0x6ba7b810, 0x9dad, 0x11d1, 0x80, 0xb4, []byte{0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}}
	URL = &UUID{0x6ba7b811, 0x9dad, 0x11d1, 0x80, 0xb4, []byte{0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}}
	OID = &UUID{0x6ba7b812, 0x9dad, 0x11d1, 0x80, 0xb4, []byte{0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}}
	DN  = &UUID{0x6ba7b814, 0x9dad, 0x11d1, 0x80, 0xb4, []byte{0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}}
)

type UUID struct {
	TimeLow   uint32
	TimeMid   uint16
	TimeHigh  uint16
	ClockHigh uint8
	ClockLow  uint8
	Node      []byte
}

func (u UUID) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, u.TimeLow)
	binary.Write(buf, binary.BigEndian, u.TimeMid)
	binary.Write(buf, binary.BigEndian, u.TimeHigh)
	binary.Write(buf, binary.BigEndian, u.ClockHigh)
	binary.Write(buf, binary.BigEndian, u.ClockLow)

	buf.Write(u.Node)

	return buf.Bytes()
}

func (u UUID) IsNil() bool {
	v := u.TimeLow == 0 && u.TimeMid == 0 && u.TimeHigh == 0 && u.ClockHigh == 0 && u.ClockLow == 0
	return v
}

func (u UUID) Version() int {
	v := (u.TimeHigh & 0xF000) >> 12
	return int(v)
}

func (u UUID) String() string {
	return fmt.Sprintf(pattern, u.TimeLow, u.TimeMid, u.TimeHigh, u.ClockHigh, u.ClockLow, u.Node)
}

func Nil() (*UUID, error) {
	return &UUID{Node: make([]byte, 6)}, nil
}

func Make(chunk []byte) (*UUID, error) {
	u := new(UUID)
	in := bytes.NewReader(chunk)
	if err := binary.Read(in, binary.BigEndian, &u.TimeLow); err != nil {
		return nil, err
	}
	if err := binary.Read(in, binary.BigEndian, &u.TimeMid); err != nil {
		return nil, err
	}
	if err := binary.Read(in, binary.BigEndian, &u.TimeHigh); err != nil {
		return nil, err
	}
	if err := binary.Read(in, binary.BigEndian, &u.ClockHigh); err != nil {
		return nil, err
	}
	if err := binary.Read(in, binary.BigEndian, &u.ClockLow); err != nil {
		return nil, err
	}
	u.Node = make([]byte, 6)
	if _, err := in.Read(u.Node); err != nil {
		return nil, err
	}
	return u, nil
}

//UUID1 create an unique identifier version 1.
func UUID1() (*UUID, error) {
	timeMU.Lock()
	defer timeMU.Unlock()

	const fraction = 100
	buf := make([]byte, 2)

	f, err := os.Open(devRandom)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	f.Read(buf)

	//epoch := time.Date(1582, 10, 15, 0, 0, 0, 0, time.Local)
	delta := time.Now().UnixNano()
	nano := uint64(delta) / fraction

	u := &UUID{
		TimeLow:   uint32(nano & 0x00000000FFFFFFFF),
		TimeMid:   uint16((nano & 0x0000FFFF00000000) >> 32),
		TimeHigh:  uint16((nano&0x0FFF000000000000)>>48) | 0x1000,
		ClockLow:  buf[1],
		ClockHigh: (buf[0] & 0x3F) | variant,
		Node:      make([]byte, 6),
	}
	is, err := net.Interfaces()
	if err != nil || len(is) == 0 {
		return nil, err
	}
	for _, ifi := range is {
		if ifi.Name != "lo" {
			u.Node = []byte(ifi.HardwareAddr)
			break
		}
	}
	return u, nil
}

//UUID1 create an unique identifier version 3 (MD5 of name is used for the different
//part of the UUID) from name and ns.
func UUID3(ns *UUID, name []byte) (*UUID, error) {
	data := append(ns.Bytes(), name...)

	sum := md5.Sum(data)
	buf := bytes.NewBuffer(sum[:])

	return read(buf, 0x3000)
}

//UUID1 create an unique identifier version 4. 
func UUID4() (*UUID, error) {
	randMU.Lock()
	defer randMU.Unlock()

	f, err := os.Open(devRandom)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return read(f, 0x4000)
}

//UUID1 create an unique identifier version 3 (SHA1 of name is used for the different
//part of the UUID) from name and ns.
func UUID5(ns *UUID, name []byte) (*UUID, error) {
	data := append(ns.Bytes(), name...)

	sum := sha1.Sum(data)
	buf := bytes.NewBuffer(sum[:])

	return read(buf, 0x5000)
}

func read(r io.Reader, version uint16) (*UUID, error) {
	u := &UUID{}
	if err := binary.Read(r, binary.BigEndian, &u.TimeLow); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &u.TimeMid); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &u.TimeHigh); err != nil {
		return nil, err
	} else {
		u.TimeHigh = version | (u.TimeHigh & 0x0FFF)
	}
	if err := binary.Read(r, binary.BigEndian, &u.ClockLow); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &u.ClockHigh); err != nil {
		return nil, err
	} else {
		u.ClockHigh |= variant
	}

	u.Node = make([]byte, 6)
	if _, err := r.Read(u.Node); err != nil {
		return nil, err
	}
	return u, nil
}
