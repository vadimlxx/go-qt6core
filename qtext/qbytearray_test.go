package qtext

import (
    "bytes"
    "testing"
)

func TestConstructors(t *testing.T) {
    e := NewQByteArray()
    defer e.Close()
    if !e.IsEmpty() || !e.IsNull() {
        t.Fatalf("New(): empty=%v null=%v", e.IsEmpty(), e.IsNull())
    }

    q := QByteArrayFromBytes([]byte("hello"))
    defer q.Close()
    if q.Size() != 5 || q.String() != "hello" {
        t.Fatalf("FromBytes: size=%d str=%q", q.Size(), q.String())
    }

    f := QByteArrayFromFill(4, 'x')
    defer f.Close()
    if f.String() != "xxxx" {
        t.Fatalf("FromFill: %q", f.String())
    }
}

func TestCOWClone(t *testing.T) {
    a := NewQByteArray("abc")
    defer a.Close()
    b := a.Clone()
    defer b.Close()

    if b.String() != "abc" {
        t.Fatalf("clone: %q", b.String())
    }
    // Mutating b must not affect a (Qt COW detach)
    b.SetAt(0, 'X')
    if a.String() != "abc" || b.String() != "Xbc" {
        t.Fatalf("COW broken: a=%q b=%q", a.String(), b.String())
    }
}

func TestElementAccess(t *testing.T) {
    q := NewQByteArray("hello")
    defer q.Close()

    if q.At(0) != 'h' || q.At(4) != 'o' {
        t.Fatalf("At: %c %c", q.At(0), q.At(4))
    }
    q.SetAt(1, 'E')
    if q.String() != "hEllo" {
        t.Fatalf("SetAt: %q", q.String())
    }
    if q.Front() != 'h' || q.Back() != 'o' {
        t.Fatalf("Front/Back: %c/%c", q.Front(), q.Back())
    }
}

func TestAppendPrependInsert(t *testing.T) {
    q := NewQByteArray("bc")
    defer q.Close()

    q.PrependString("a")
    q.AppendString("d")
    q.InsertString(2, "X")
    if q.String() != "abXcd" {
        t.Fatalf("got %q, want %q", q.String(), "abXcd")
    }

    q.Remove(2, 1)
    if q.String() != "abcd" {
        t.Fatalf("after Remove: %q", q.String())
    }
    q.RemoveFirst()
    q.RemoveLast()
    if q.String() != "bc" {
        t.Fatalf("after RemoveFirst/Last: %q", q.String())
    }
}

func TestSearchAndSlice(t *testing.T) {
    q := NewQByteArray("the quick brown fox")
    defer q.Close()

    if q.IndexOf('q', 0) != 4 {
        t.Fatalf("IndexOf('q'): %d", q.IndexOf('q', 0))
    }
    if q.LastIndexOf('o', -1) != 17 {
        t.Fatalf("LastIndexOf('o'): %d", q.LastIndexOf('o', -1))
    }
    if !q.ContainsBytes([]byte("brown")) {
        t.Fatal("ContainsBytes(brown)")
    }
    if q.Count(' ') != 3 {
        t.Fatalf("Count(' '): %d", q.Count(' '))
    }
    if !q.StartsWithBytes([]byte("the")) || !q.EndsWithBytes([]byte("fox")) {
        t.Fatal("StartsWith/EndsWith")
    }

    l := q.Left(3)
    defer l.Close()
    if l.String() != "the" {
        t.Fatalf("Left: %q", l.String())
    }
    r := q.Right(3)
    defer r.Close()
    if r.String() != "fox" {
        t.Fatalf("Right: %q", r.String())
    }
    m := q.Mid(4, 5)
    defer m.Close()
    if m.String() != "quick" {
        t.Fatalf("Mid: %q", m.String())
    }
}

func TestTransforms(t *testing.T) {
    q := NewQByteArray("  Hello, World!  ")
    defer q.Close()

    lo := q.ToLower()
    defer lo.Close()
    if lo.String() != "  hello, world!  " {
        t.Fatalf("ToLower: %q", lo.String())
    }

    tr := q.Trimmed()
    defer tr.Close()
    if tr.String() != "Hello, World!" {
        t.Fatalf("Trimmed: %q", tr.String())
    }

    rep := NewQByteArray("ab").Repeated(3)
    defer rep.Close()
    if rep.String() != "ababab" {
        t.Fatalf("Repeated: %q", rep.String())
    }
}

func TestNumeric(t *testing.T) {
    q := NewQByteArray("42")
    defer q.Close()

    v, ok := q.ToInt(10)
    if !ok || v != 42 {
        t.Fatalf("ToInt: %d %v", v, ok)
    }

    bad := NewQByteArray("not-a-number")
    defer bad.Close()
    if _, ok := bad.ToInt(10); ok {
        t.Fatal("ToInt should fail")
    }

    n := NumberLL(-123, 10)
    defer n.Close()
    if n.String() != "-123" {
        t.Fatalf("NumberLL: %q", n.String())
    }
}

func TestBase64(t *testing.T) {
    src := NewQByteArray("hello world")
    defer src.Close()

    b64 := src.ToBase64(0)
    defer b64.Close()
    if b64.String() != "aGVsbG8gd29ybGQ=" {
        t.Fatalf("ToBase64: %q", b64.String())
    }

    dec, status := FromBase64Encoding([]byte("aGVsbG8gd29ybGQ="), 0)
    defer dec.Close()
    if status != Base64DecodingStatusOk || dec.String() != "hello world" {
        t.Fatalf("FromBase64: %q status=%v", dec.String(), status)
    }

    bad, status := FromBase64Encoding([]byte("!!!not-base64!!!"), Base64Options(AbortOnBase64DecodingErrors))
    if bad != nil {
        bad.Close()
    }
    if status == Base64DecodingStatusOk {
        t.Fatal("expected non-Ok status")
    }
}

func TestHex(t *testing.T) {
    src := NewQByteArray([]byte{0xDE, 0xAD, 0xBE, 0xEF})
    defer src.Close()

    h := src.ToHex(0)
    defer h.Close()
    if h.String() != "deadbeef" {
        t.Fatalf("ToHex: %q", h.String())
    }

    back := FromHex([]byte("deadbeef"))
    defer back.Close()
    if !bytes.Equal(back.Bytes(), []byte{0xDE, 0xAD, 0xBE, 0xEF}) {
        t.Fatalf("FromHex: %v", back.Bytes())
    }

    sep := src.ToHex(':')
    defer sep.Close()
    if sep.String() != "de:ad:be:ef" {
        t.Fatalf("ToHex(':'): %q", sep.String())
    }
}

func TestPercent(t *testing.T) {
    src := NewQByteArray("a b/c?d")
    defer src.Close()

    p := src.ToPercentEncoding(nil, nil, '%')
    defer p.Close()
    if p.String() != "a%20b%2Fc%3Fd" {
        t.Fatalf("ToPercent: %q", p.String())
    }

    back := FromPercentEncoding([]byte("a%20b%2Fc%3Fd"), '%')
    defer back.Close()
    if back.String() != "a b/c?d" {
        t.Fatalf("FromPercent: %q", back.String())
    }
}

func TestCompress(t *testing.T) {
    src := NewQByteArray(string(bytes.Repeat([]byte("hello "), 100)))
    defer src.Close()

    c := src.Compress(-1)
    defer c.Close()
    if c.Size() >= src.Size() {
        t.Fatalf("compressed should be smaller: %d vs %d", c.Size(), src.Size())
    }

    u := c.Uncompress()
    defer u.Close()
    if u.String() != src.String() {
        t.Fatal("compress/uncompress round-trip failed")
    }
}

func TestSplit(t *testing.T) {
    src := NewQByteArray("a,b,c,,d")
    defer src.Close()

    parts := src.Split(',')
    if len(parts) != 5 {
        t.Fatalf("Split: %d parts", len(parts))
    }
    for _, p := range parts {
        defer p.Close()
    }
    want := []string{"a", "b", "c", "", "d"}
    for i, p := range parts {
        if p.String() != want[i] {
            t.Fatalf("part[%d] = %q, want %q", i, p.String(), want[i])
        }
    }
}

func TestUtf8(t *testing.T) {
    ok := NewQByteArray("héllo")
    defer ok.Close()
    if !ok.IsValidUtf8() {
        t.Fatal("valid utf8 expected")
    }

    bad := NewQByteArray([]byte{0xFF, 0xFE, 0xFD})
    defer bad.Close()
    if bad.IsValidUtf8() {
        t.Fatal("invalid utf8 expected")
    }
}
