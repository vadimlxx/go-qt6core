package quuid

import (
	"testing"

	"github.com/vadimlxx/go-qt6core/qdata"
)

func TestQUuidNull(t *testing.T) {
	q := NewQUuid()
	defer q.Close()
	if !q.IsNull() {
		t.Fatal("NewQUuid should be null")
	}
	if q.Variant() != VarUnknown {
		t.Fatalf("null variant: %d", q.Variant())
	}
	if q.Version() != VerUnknown {
		t.Fatalf("null version: %d", q.Version())
	}
}

func TestQUuidFromUint32(t *testing.T) {
	q := NewQUuidFromUint32(0x12345678, 0x9abc, 0xdef0, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0)
	defer q.Close()
	if q.IsNull() {
		t.Fatal("expected non-null")
	}
	got := q.ToString(WithBraces)
	if got != "{12345678-9abc-def0-1234-56789abcdef0}" {
		t.Fatalf("ToString: %q", got)
	}
}

func TestQUuidStringRoundTrip(t *testing.T) {
	const src = "12345678-9abc-def0-1234-56789abcdef0"
	q := NewQUuidFromString(src)
	defer q.Close()
	if q.IsNull() {
		t.Fatal("parse failed")
	}
	got := q.ToString(WithoutBraces)
	if got != src {
		t.Fatalf("round-trip: %q", got)
	}

	got2 := q.ToString(WithBraces)
	if got2 != "{"+src+"}" {
		t.Fatalf("WithBraces: %q", got2)
	}

	got3 := q.ToString(Id128)
	if got3 != "123456789abcdef0123456789abcdef0" {
		t.Fatalf("Id128: %q", got3)
	}
}

func TestQUuidToByteArray(t *testing.T) {
	q := NewQUuidFromString("12345678-9abc-def0-1234-56789abcdef0")
	defer q.Close()
	b := q.ToByteArray(WithoutBraces)
	if string(b) != "12345678-9abc-def0-1234-56789abcdef0" {
		t.Fatalf("ToByteArray: %q", b)
	}
}

func TestQUuidBytesRoundTrip(t *testing.T) {
	var raw Id128Bytes
	for i := range raw {
		raw[i] = byte(i)
	}
	q := NewQUuidFromBytes(raw, qdata.BigEndian)
	defer q.Close()
	back := q.ToBytes(qdata.BigEndian)
	if back != raw {
		t.Fatalf("BigEndian round-trip: %v != %v", back, raw)
	}

	q2 := NewQUuidFromBytes(raw, qdata.LittleEndian)
	defer q2.Close()
	back2 := q2.ToBytes(qdata.LittleEndian)
	if back2 != raw {
		t.Fatalf("LittleEndian round-trip: %v != %v", back2, raw)
	}
}

func TestQUuidRfc4122(t *testing.T) {
	const src = "12345678-9abc-def0-1234-56789abcdef0"
	q := NewQUuidFromString(src)
	defer q.Close()
	b := q.ToRfc4122()
	if len(b) != 16 {
		t.Fatalf("Rfc4122 len: %d", len(b))
	}
	q2 := NewQUuidFromRfc4122(b)
	defer q2.Close()
	if !q.Equal(q2) {
		t.Fatal("Rfc4122 round-trip mismatch")
	}
}

func TestQUuidCreateUuid(t *testing.T) {
	q1 := CreateUuid()
	q2 := CreateUuid()
	defer q1.Close()
	defer q2.Close()
	if q1.IsNull() || q2.IsNull() {
		t.Fatal("CreateUuid returned null")
	}
	if q1.Equal(q2) {
		t.Fatal("two random UUIDs should differ")
	}
	if q1.Version() != Random {
		t.Fatalf("version: %d", q1.Version())
	}
}

func TestQUuidCreateUuidV3V5(t *testing.T) {
	ns := NewQUuidFromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	defer ns.Close()

	v3 := CreateUuidV3(ns, []byte("hello"))
	defer v3.Close()
	if v3.Version() != Md5 {
		t.Fatalf("v3 version: %d", v3.Version())
	}

	v5 := CreateUuidV5(ns, []byte("hello"))
	defer v5.Close()
	if v5.Version() != Sha1 {
		t.Fatalf("v5 version: %d", v5.Version())
	}

	// Same input => same output (deterministic).
	v3b := CreateUuidV3(ns, []byte("hello"))
	defer v3b.Close()
	if !v3.Equal(v3b) {
		t.Fatal("v3 not deterministic")
	}
}

func TestQUuidCreateUuidV7(t *testing.T) {
	q := CreateUuidV7()
	defer q.Close()
	if q.IsNull() {
		t.Fatal("CreateUuidV7 returned null")
	}
	if q.Version() != UnixEpoch {
		t.Fatalf("v7 version: %d", q.Version())
	}
}

func TestQUuidCompare(t *testing.T) {
	a := NewQUuidFromString("00000000-0000-0000-0000-000000000001")
	b := NewQUuidFromString("00000000-0000-0000-0000-000000000002")
	defer a.Close()
	defer b.Close()
	if a.Compare(b) >= 0 {
		t.Fatalf("Compare: %d", a.Compare(b))
	}
	if b.Compare(a) <= 0 {
		t.Fatalf("Compare: %d", b.Compare(a))
	}
	if a.Compare(a) != 0 {
		t.Fatalf("self Compare: %d", a.Compare(a))
	}
}