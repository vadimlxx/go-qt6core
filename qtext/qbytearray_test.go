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

    q := NewQByteArrayFromBytes([]byte("hello"))
    defer q.Close()
    if q.Size() != 5 || q.String() != "hello" {
        t.Fatalf("FromBytes: size=%d str=%q", q.Size(), q.String())
    }

    f := NewQByteArrayFromFill(4, 'x')
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

    n := NewQByteArrayNumberLL(-123, 10)
    defer n.Close()
    if n.String() != "-123" {
        t.Fatalf("NewQByteArrayNumberLL: %q", n.String())
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

    dec, status := NewQByteArrayFromBase64Encoding([]byte("aGVsbG8gd29ybGQ="), 0)
    defer dec.Close()
    if status != Base64DecodingStatusOk || dec.String() != "hello world" {
        t.Fatalf("FromBase64: %q status=%v", dec.String(), status)
    }

    bad, status := NewQByteArrayFromBase64Encoding([]byte("!!!not-base64!!!"), Base64Options(AbortOnBase64DecodingErrors))
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

    back := NewQByteArrayFromHex([]byte("deadbeef"))
    defer back.Close()
    if !bytes.Equal(back.Bytes(), []byte{0xDE, 0xAD, 0xBE, 0xEF}) {
        t.Fatalf("NewQByteArrayFromHex: %v", back.Bytes())
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

    back := NewQByteArrayFromPercentEncoding([]byte("a%20b%2Fc%3Fd"), '%')
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

func TestAssignAndSwap(t *testing.T) {
    a := NewQByteArray("left")
    defer a.Close()
    b := NewQByteArray("right")
    defer b.Close()

    a.Swap(b)
    if a.String() != "right" || b.String() != "left" {
        t.Fatalf("Swap: a=%q b=%q", a.String(), b.String())
    }

    a.AssignBytes([]byte("abc"))
    if a.String() != "abc" {
        t.Fatalf("AssignBytes: %q", a.String())
    }
    a.AssignString("xyz")
    if a.String() != "xyz" {
        t.Fatalf("AssignString: %q", a.String())
    }
    a.AssignFill(3, 'q')
    if a.String() != "qqq" {
        t.Fatalf("AssignFill: %q", a.String())
    }
    src := NewQByteArray("copy")
    defer src.Close()
    a.AssignQByteArray(src)
    if a.String() != "copy" {
        t.Fatalf("AssignQByteArray: %q", a.String())
    }
}

func TestNullEmptySemantics(t *testing.T) {
    n := NewQByteArrayFromBytes(nil)
    defer n.Close()
    if !n.IsNull() || !n.IsEmpty() {
        t.Fatalf("nil bytes: null=%v empty=%v", n.IsNull(), n.IsEmpty())
    }

    e := NewQByteArrayFromBytes([]byte{})
    defer e.Close()
    if e.IsNull() || !e.IsEmpty() {
        t.Fatalf("empty bytes: null=%v empty=%v", e.IsNull(), e.IsEmpty())
    }

    e.AssignBytes([]byte("x"))
    e.AssignBytes([]byte{})
    if e.IsNull() || !e.IsEmpty() {
        t.Fatalf("AssignBytes(empty): null=%v empty=%v", e.IsNull(), e.IsEmpty())
    }
    e.AssignBytes(nil)
    if !e.IsNull() {
        t.Fatal("AssignBytes(nil) should make the array null")
    }
}

func TestSlice(t *testing.T) {
    q := NewQByteArray("abcdef")
    defer q.Close()
    q.Slice(2)
    if q.String() != "cdef" {
        t.Fatalf("Slice: %q", q.String())
    }

    q2 := NewQByteArray("abcdef")
    defer q2.Close()
    q2.SliceN(1, 3)
    if q2.String() != "bcd" {
        t.Fatalf("SliceN: %q", q2.String())
    }
}

func TestInsertReplaceVariants(t *testing.T) {
    q := NewQByteArray("ac")
    defer q.Close()
    ins := NewQByteArray("b")
    defer ins.Close()
    q.InsertQByteArray(1, ins)
    if q.String() != "abc" {
        t.Fatalf("InsertQByteArray: %q", q.String())
    }

    q.ReplaceCharBytes('b', []byte("XY"))
    if q.String() != "aXYc" {
        t.Fatalf("ReplaceCharBytes: %q", q.String())
    }
    q.ReplaceAllString("XY", "b")
    if q.String() != "abc" {
        t.Fatalf("ReplaceAllString: %q", q.String())
    }
}

func TestPercentDecoded(t *testing.T) {
    q := NewQByteArray("a%20b%2Fc")
    defer q.Close()

    d := q.PercentDecoded('%')
    defer d.Close()
    if d.String() != "a b/c" {
        t.Fatalf("PercentDecoded: %q", d.String())
    }
    if got := string(q.PercentDecodedBytes('%')); got != "a b/c" {
        t.Fatalf("PercentDecodedBytes: %q", got)
    }
}

func TestNumericVariants(t *testing.T) {
    q := NewQByteArray("123")
    defer q.Close()
    if v, ok := q.ToLong(10); !ok || v != 123 {
        t.Fatalf("ToLong: %d %v", v, ok)
    }
    if v, ok := q.ToULong(10); !ok || v != 123 {
        t.Fatalf("ToULong: %d %v", v, ok)
    }

    n := NewQByteArrayNumberInt(-42, 10)
    defer n.Close()
    if n.String() != "-42" {
        t.Fatalf("NumberInt: %q", n.String())
    }
    u := NewQByteArrayNumberUInt(42, 10)
    defer u.Close()
    if u.String() != "42" {
        t.Fatalf("NumberUInt: %q", u.String())
    }

    q.SetNumInt(-7, 10)
    if q.String() != "-7" {
        t.Fatalf("SetNumInt: %q", q.String())
    }
    q.SetNumUInt(7, 10)
    if q.String() != "7" {
        t.Fatalf("SetNumUInt: %q", q.String())
    }
    q.SetNumFloat(1.5, 'f', 2)
    if q.String() != "1.50" {
        t.Fatalf("SetNumFloat: %q", q.String())
    }
}

func TestSetFrontBack(t *testing.T) {
    q := NewQByteArray("abc")
    defer q.Close()
    q.SetFront('X')
    q.SetBack('Z')
    if q.String() != "XbZ" {
        t.Fatalf("SetFront/SetBack: %q", q.String())
    }
}

func TestRemoveIf(t *testing.T) {
    q := NewQByteArray("a1b2c3")
    defer q.Close()
    removed := q.RemoveIf(func(b byte) bool { return b >= '0' && b <= '9' })
    if removed != 3 || q.String() != "abc" {
        t.Fatalf("RemoveIf: removed=%d got=%q", removed, q.String())
    }
}

func TestConcat(t *testing.T) {
    a := NewQByteArray("foo")
    defer a.Close()
    b := NewQByteArray("bar")
    defer b.Close()

    c := a.Concat(b)
    defer c.Close()
    if c.String() != "foobar" || a.String() != "foo" {
        t.Fatalf("Concat: %q (a=%q)", c.String(), a.String())
    }
    if got := string(a.ConcatBytes([]byte("!"))); got != "foo!" {
        t.Fatalf("ConcatBytes: %q", got)
    }
}

func TestBytesVariants(t *testing.T) {
    q := NewQByteArray("  Hi  ")
    defer q.Close()

    if got := string(q.LeftBytes(4)); got != "  Hi" {
        t.Fatalf("LeftBytes: %q", got)
    }
    if got := string(q.ToUpperBytes()); got != "  HI  " {
        t.Fatalf("ToUpperBytes: %q", got)
    }
    if got := string(q.TrimmedBytes()); got != "Hi" {
        t.Fatalf("TrimmedBytes: %q", got)
    }
    if got := string(q.CloneBytes()); got != "  Hi  " {
        t.Fatalf("CloneBytes: %q", got)
    }

    hexSrc := NewQByteArray([]byte{0xAB, 0xCD})
    defer hexSrc.Close()
    if got := string(hexSrc.ToHexBytes(0)); got != "abcd" {
        t.Fatalf("ToHexBytes: %q", got)
    }
    if got := string(hexSrc.ToBase64Bytes(0)); got != "q80=" {
        t.Fatalf("ToBase64Bytes: %q", got)
    }

    nt := NewQByteArray("x")
    defer nt.Close()
    if got := nt.NullTerminatedBytes(); len(got) != 1 || got[0] != 'x' {
        t.Fatalf("NullTerminatedBytes: %v", got)
    }
}

func TestSplitBytes(t *testing.T) {
    q := NewQByteArray("a,b,,c")
    defer q.Close()
    parts := q.SplitBytes(',')
    if len(parts) != 4 {
        t.Fatalf("SplitBytes: %d parts", len(parts))
    }
    want := []string{"a", "b", "", "c"}
    for i, p := range parts {
        if string(p) != want[i] {
            t.Fatalf("part[%d] = %q, want %q", i, p, want[i])
        }
    }
}

func TestNewFromString(t *testing.T) {
    q := NewQByteArrayFromString("héllo\x00world")
    defer q.Close()
    if q.Size() != 12 {
        t.Fatalf("NewQByteArrayFromString: size=%d", q.Size())
    }
    if q.At(1) != 0xC3 || q.At(2) != 0xA9 {
        t.Fatalf("UTF-8 bytes expected, got %x %x", q.At(1), q.At(2))
    }
}
