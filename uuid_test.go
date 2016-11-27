package uuid

import (
	"bytes"
	"testing"
)

func TestNilUUID(t *testing.T) {
	uid, _ := Nil()
	if !uid.IsNil() {
		t.Error("not nil uuid")
	}
}

func TestUUID1Uniqueness(t *testing.T) {
	seens := make(map[*UUID]struct{})
	for i := 0; i < 10000; i++ {
		uid, _ := UUID1()
		if _, ok := seens[uid]; ok {
			t.Errorf("%s not unique", uid)
		}
		seens[uid] = struct{}{}
	}
}

func TestUUID4Uniqueness(t *testing.T) {
	seens := make(map[*UUID]struct{})
	for i := 0; i < 10000; i++ {
		uid, _ := UUID4()
		if _, ok := seens[uid]; ok {
			t.Errorf("%s not unique", uid)
		}
		seens[uid] = struct{}{}
	}
}

func TestUUIDVersionAndLen(t *testing.T) {
	type pair struct {
		f    func() (*UUID, error)
		want int
	}
	data := []pair{
		{UUID1, 1},
		{UUID4, 4},
	}

	for _, d := range data {
		uid, err := d.f()
		if err != nil {
			t.Error(err)
		}
		if v := uid.Version(); v != d.want {
			t.Errorf("wrong version! expected: %d, got: %d", d.want, v)
		}
		if b := uid.Bytes(); len(b) != 16 {
			t.Errorf("wrong length! expected: 16, got: %d", len(b))
		}
		if uid.IsNil() {
			t.Errorf("uuid is nil uuid")
		}
	}
}

func TestUUIDNameVersionAndLen(t *testing.T) {
	type pair struct {
		f    func(*UUID, []byte) (*UUID, error)
		want int
	}
	data := []pair{
		{UUID3, 3},
		{UUID5, 5},
	}

	for _, d := range data {
		uid, err := d.f(DNS, []byte("hello"))
		if err != nil {
			t.Error(err)
			continue
		}
		if v := uid.Version(); v != d.want {
			t.Errorf("wrong version! expected: %d, got: %d", d.want, v)
		}
		if b := uid.Bytes(); len(b) != 16 {
			t.Errorf("wrong length! expected: 16, got: %d", len(b))
		}
		if uid.IsNil() {
			t.Errorf("uuid is nil uuid")
		}
	}
}

func TestUUID5Equality(t *testing.T) {
	var uid1, uid2 *UUID

	uid1, _ = UUID5(DNS, []byte("hello"))
	uid2, _ = UUID5(DNS, []byte("hello"))
	if !bytes.Equal(uid1.Bytes(), uid2.Bytes()) {
		t.Errorf("%s and %s not equal: expected equal", uid1, uid2)
	}

	uid1, _ = UUID5(DNS, []byte("hello"))
	uid2, _ = UUID5(DNS, []byte("world"))
	if bytes.Equal(uid1.Bytes(), uid2.Bytes()) {
		t.Errorf("%s and %s equal: expected not equal", uid1, uid2)
	}

	uid1, _ = UUID5(DNS, []byte("hello"))
	uid2, _ = UUID5(URL, []byte("hello"))
	if bytes.Equal(uid1.Bytes(), uid2.Bytes()) {
		t.Errorf("%s and %s equal: expected not equal", uid1, uid2)
	}
}

func TestUUID3Equality(t *testing.T) {
	var uid1, uid2 *UUID

	uid1, _ = UUID3(DNS, []byte("hello"))
	uid2, _ = UUID3(DNS, []byte("hello"))
	if !bytes.Equal(uid1.Bytes(), uid2.Bytes()) {
		t.Errorf("%s and %s not equal: expected equal", uid1, uid2)
	}

	uid1, _ = UUID3(DNS, []byte("hello"))
	uid2, _ = UUID3(DNS, []byte("world"))
	if bytes.Equal(uid1.Bytes(), uid2.Bytes()) {
		t.Errorf("%s and %s equal: expected not equal", uid1, uid2)
	}

	uid1, _ = UUID3(DNS, []byte("hello"))
	uid2, _ = UUID3(URL, []byte("hello"))
	if bytes.Equal(uid1.Bytes(), uid2.Bytes()) {
		t.Errorf("%s and %s equal: expected not equal", uid1, uid2)
	}
}
