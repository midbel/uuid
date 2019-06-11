package uuid

import (
	"io"
	"fmt"
	"crypto/md5"
	"crypto/sha1"
	"crypto/rand"
)

const Size = 16

const variant byte = 0x40

const pattern =  "%08x-%04x-%04x-%02x%02x-%x"

var (
	// 6ba7b810-9dad-11d1-80b4-00c04fd430c8
	DNS = UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	// 6ba7b811-9dad-11d1-80b4-00c04fd430c8
	URL = UUID{0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	// 	6ba7b812-9dad-11d1-80b4-00c04fd430c8
	OID = UUID{0x6b, 0xa7, 0xb8, 0x12, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	// 6ba7b814-9dad-11d1-80b4-00c04fd430c8
	DN  = UUID{0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
)

var Nil = UUID(make([]byte, Size))

type UUID []byte

func (u UUID) String() string {
	xs := []byte(u)

	var (
		tlow  = xs[:4]
		tmid  = xs[4:6]
		thigh = xs[6:8]
		chigh = xs[8]
		clow  = xs[9]
		node  = xs[10:]
	)

	return fmt.Sprintf(pattern, tlow, tmid, thigh, chigh, clow, node)
}

func (u UUID) Version() int {
	return int(u[6]) >> 4
}

// func UUID1() UUID {
// 	return nil
// }

func UUID4() UUID {
	xs := make([]byte, Size)
	if _, err := io.ReadFull(rand.Reader, xs); err != nil {
		return Nil
	}
	return update(xs, 0x40)
}

func UUID3(str []byte, ns UUID) UUID {
	sum := md5.Sum(append([]byte(ns), str...))
	return update(sum[:], 0x30)
}

func UUID5(str []byte, ns UUID) UUID {
	sum := sha1.Sum(append([]byte(ns), str...))
	return update(sum[:], 0x50)
}

func update(xs []byte, version byte) UUID {
	if len(xs) > Size {
		xs = xs[:Size]
	}
	xs[6] = (xs[6] & 0x0F) | version
	xs[8] += variant

	return UUID(xs)
}
